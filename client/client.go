package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func main() {
	url := flag.String("url", "http://localhost:3000", "the URL of the Newmarket server")
	port := flag.Int("port", 3001, "The port you want to listen on")
	bindlocal := flag.Bool("bindlocal", true, "enable to bind only on 127.0.0.1")
	flag.Parse()
	StartTunnel(*url, int64(*port), *bindlocal)
}

func StartTunnel(URL string, Port int64, BindLocal bool) {
	fmt.Printf("The settings are \n\nURL:%s\nListening Port:%s\n", URL, Port)
	// First, Lets see if we can bind that port.
	var bindaddr string = "0.0.0.0"

	if BindLocal {
		bindaddr = "127.0.0.1"
	}

	listener, e := net.Listen("tcp", fmt.Sprintf("%s:%d", bindaddr, Port))
	if e != nil {
		fmt.Errorf("Cannot bind to port %s:%d", bindaddr, Port)
		return
	}
	fmt.Printf("Bound to %s:%d waiting for a connection to proceed\n", bindaddr, Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("Error accept incoming connection: %s", err.Error())
			return
		}
		go HandleTunConnection(conn, URL, Port)
	}

}

func HandleTunConnection(conn net.Conn, URL string, Port int64) {
	fmt.Println("Getting a session ID from server...\n")
	r, e := http.Get(URL + "/init")
	if e != nil {
		fmt.Errorf("Unable to get a session!")
		conn.Close()
		return
	}
	sessdata, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Errorf("Unable to read from the connection to get a session!")
		conn.Close()
		return
	}
	sessiontoken := string(sessdata)
	fmt.Println("Session tokean obtained:", sessiontoken)
	// Okay so we now have our session token.
	go DialUpWards(URL, sessiontoken, conn)
	go DialDownWards(URL+"/session/"+sessiontoken, conn)
}

func DialUpWards(URL string, sessiontoken string, conn net.Conn) {
	// because go can't do what I am about to do I am going to
	// #yolo my own HTTP code fora bit :v
	URL = strings.Replace(URL, "http://", "", 1)
	// WARNING THIS CODE WONT HELP YOU AGAINST SUB FOLDER URLS
	// NEED TO BE REFACTORED OR HELL, JUST DONE PROPERALLY AND NOT
	// INSANE LIKE THIS ONE IS.
	conn2, err := net.Dial("tcp", URL)
	defer conn.Close()
	defer conn2.Close()
	if err != nil {
		fmt.Errorf("Woah wtf, I tried to dial up and I got a error! %s", err)
		return
	}
	HTTPRequest := fmt.Sprintf("POST /session/%s HTTP/1.1\r\nHost: %s\r\nUser-Agent: Newmarket\r\nContent-Length: 9999999\r\n\r\n", sessiontoken, URL)
	conn2.Write([]byte(HTTPRequest))
	buf := make([]byte, 25565)
	for {
		read, e := conn.Read(buf)
		if e != nil {
			fmt.Errorf("Upstream broke down for reason %s", e.Error())
			return
		}
		_, e = conn2.Write(buf[0:read])
		if e != nil {
			fmt.Errorf("Tried to write data to remotesocket and it broke: %s", e.Error())
			return
		}
	}
}

func DialDownWards(URL string, conn net.Conn) {
	r, e := http.Get(URL)
	defer conn.close()
	if e != nil {
		fmt.Errorf("Woah wtf, I tried to dial down and I got a error! %s", e.Error())
		return
	}
	buf := make([]byte, 25565)
	for {
		read, e := r.Body.Read(buf)
		if e != nil {
			fmt.Errorf("Downstream broke down for reason %s", e.Error())
			return
		}
		_, e = conn.Write(buf[0:read])
		if e != nil {
			fmt.Errorf("Tried to write data to localsocket and it broke: %s", e.Error())
			return
		}
	}
}
