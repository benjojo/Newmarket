package main

import (
	// "fmt"
	"fmt"
	"github.com/codegangsta/martini"
	// "io"
	"net"
	"net/http"
)

func main() {
	UpChan := make(chan []byte)
	DoChan := make(chan []byte)
	go TCPSocket(UpChan, DoChan)
	// Now that the TCP waiter is setup. lets start the HTTP sevrer
	m := martini.Classic()
	// m.Use(EnforceHTTPAuth)
	m.Get("/", Welcome)
	m.Run()
}

func TCPSocket(UpChan chan []byte, DoChan chan []byte) {
	// To start off the connection we expect there to be a "ping"
	// on the UpChan.
	<-UpChan
	// This blocks the first "contact"
	// and awakes the server up from its terrifying slumber
	_, err := net.Dial("tcp", "localhost:22")
	if err != nil {
		fmt.Errorf("Could not dial SSH on the localhost, this is a srs issue. %s", err)
	}
	/*
		Okay so you need to first do somthing with that _ up there.
		I think it will be best to move this to use one chan and just have a struct that can contain stuff
		to move back and forward, that way everything can be streamlined into a multi connection system
		thus making the system much more sane a reliable.
	*/
}

func Welcome(rw http.ResponseWriter, req *http.Request) string {
	return "Why Howdy there"
}
