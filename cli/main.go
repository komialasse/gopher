package main

import (
	"flag"
	"github.com/gopher"
)

func main() {
	var host string
	var to string
	flag.StringVar(&host, "host", "localhost", "the host to expose")
	flag.StringVar(&to, "to", "server", "address of the remote server")
	flag.Parse()


	server := gopher.Server{}
	server.Listen()
}
