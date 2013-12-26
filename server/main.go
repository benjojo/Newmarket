package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/nu7hatch/gouuid"
	// "io"
	"net"
	"net/http"
)

type ConnectionSession struct {
	Token    string
	UpChan   chan []byte
	DownChan chan []byte
}

var Sessions = make([]ConnectionSession, 0)

func main() {
	// Now that the TCP waiter is setup. lets start the HTTP sevrer
	m := martini.Classic()
	// m.Map(Sessions)
	m.Get("/", Welcome)
	m.Get("/init", StartSession)
	/*
		Okay so all I need to do now is make a m.get(/uplink/:id) and a post one too to get the data
		You can apparently stream r.body. you can flush to the output as long as you DONT use martini's
		fancy way of doing things and just write to the request in raw form.
	*/
	m.Run()
}

func StartSession(rw http.ResponseWriter, req *http.Request) string {
	// Now we need to make a new session and store it in a KV DB
	UpChan := make(chan []byte)
	DownChan := make(chan []byte)
	u, _ := uuid.NewV4()
	ustr := fmt.Sprintf("%s", u)
	WorkingSession := ConnectionSession{
		Token:    ustr,
		UpChan:   UpChan,
		DownChan: DownChan,
	}
	go TCPSocket(WorkingSession)
	Sessions = append(Sessions, WorkingSession)
	return ustr
}

func UpPoll(conn net.Conn, UpChan chan []byte) {
	for {
		conn.Write(<-UpChan)
	}
}

func DownPoll(conn net.Conn, DownChan chan []byte) {
	for {
		buf := make([]byte, 25565)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Errorf("Could not Read!!! %s", err)
			break
		} else {
			DownChan <- buf[:n]
		}
	}
}

func TCPSocket(Session ConnectionSession) {
	<-Session.UpChan
	UpChan := Session.UpChan
	DownChan := Session.DownChan

	// This blocks the first "contact"
	// and awakes the server up from its terrifying slumber
	conn, err := net.Dial("tcp", "localhost:22")
	if err != nil {
		fmt.Errorf("Could not dial SSH on the localhost, this is a srs issue. %s", err)
	}
	go UpPoll(conn, UpChan)
	go DownPoll(conn, DownChan)
}

func Welcome(rw http.ResponseWriter, req *http.Request) string {
	return "Why Howdy there"
}
