package main

import (
	// "fmt"
	"github.com/codegangsta/martini"
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

}

func Welcome(rw http.ResponseWriter, req *http.Request) string {
	return "Why Howdy there"
}
