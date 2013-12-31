package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "Newmarket Client"
	app.Usage = "A Client to the Newmarket HTTP Tunnel server"
	app.Action = func(c *cli.Context) {
		fmt.Println("Starting Newmarket client")
		StartTunnel(c.String("url"), c.String("port"))
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{"url, u", "http://localhost:3000", "the URL of the Newmarket server"},
		cli.StringFlag{"port, p", "3001", "The port you want to listen on"},
	}
	app.Run(os.Args)
}

func StartTunnel(URL string, Port string) {
	fmt.Printf("The settings are \n\nURL:%s\nListening Port:%s\n", URL, Port)
	// First, Lets see if we can bind that port.
	i, e := strconv.ParseInt(Port, 10, 64)
	if e != nil {
		fmt.Errorf("The port '%s' is not a valid int. wtf did you put in?!", Port)
		return
	}
}
