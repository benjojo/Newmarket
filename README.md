Newmarket
=========

Wrap a TCP connection over two HTTP connections like so:

![](http://i.imgur.com/V1kKCb1.png)

To a firewall the two connections looks like two (correctly formatted) HTTP connections that are slowly POSTing and GETing data from a HTTP url.

##Setup
Make sure you have GoLang 1.1 or above installed then clone the git.

###Server
First go into the server directory
`cd server`

Fetch the things that are needed to run this first:

`go get`

then build it

`go build`

You can then run it, If you want to change the port that it listens on, set the PORT env var by doing (in a normal bash shell) `export port=80`

then run the server program

`./server`

###Client

First go into the server directory
`cd client`

Fetch the things that are needed to run this first:

`go get`

then build it

`go build`

Then you can look into the usage of the client by running `./client --help` or `client.exe --help`

You will get this in return

```
NAME:
   Newmarket Client - A Client to the Newmarket HTTP Tunnel server

USAGE:
   Newmarket Client [global options] command [command options] [arguments...]

VERSION:
   0.9

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url, -u 'http://localhost:3000'    the URL of the Newmarket server
   --port, -p '3001'                    The port you want to listen on
   --version, -v                        print the version
   --help, -h                           show help
```

The two things you will want to change here is the `--url` var to point to your remote server that you will be tunneling to.

Example usage is

`client --url http://test.example.com`

then connect on the forwarded port (default is 3001)

`ssh localhost -P 3001`

Auth like normal and now you are inside your remote server!
