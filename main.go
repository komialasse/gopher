package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gopher"
)

const (
	LOCAL  = "local"
	SERVER = "server"
)



func main() {

	gopher.RegisterMessages()

	command := os.Args[1]

	switch command {
	case LOCAL:
		localPort, err := strconv.Atoi(os.Args[2])
		subArgs := os.Args[3:]
		if err != nil {
			log.Fatalf("Could not parse local port %v", os.Args[2])
		}

		localHost, to, remotePort := parseLocalArgs(subArgs)
		client := gopher.NewClient(localHost, to, localPort, remotePort)
		client.Listen(context.Background())
	case SERVER:
		server := gopher.NewServer()
		server.Listen(context.Background())
	default:
		log.Fatal("unrecognized command")
	}
}

func parseLocalArgs(args []string) (string, string, int) {
	localfs := flag.NewFlagSet(LOCAL, flag.ExitOnError)
	var localhost string
	var port int 
	localfs.StringVar(&localhost, "localhost", "localhost", "the local host to expose")
	localfs.StringVar(&localhost, "l", "localhost", "the local host to expose")
	to := localfs.String("to", "localhost", "address of remote server")
	localfs.IntVar(&port, "port", 8081, "port of remote server")
	localfs.IntVar(&port, "p", 8081, "port of remote server")

	if err := localfs.Parse(args); err != nil {
		log.Fatal(err)
	}
	return localhost, *to, port
}
