package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gopher"
)

const (
	LOCAL = "local"
	SERVER = "server"
)

func main() {
	command := os.Args[1]

	switch command {
	case LOCAL:
		localPort, err := strconv.Atoi(os.Args[2])
		subArgs := os.Args[3:]
		if err != nil {
			log.Fatalf("Could not parse local port %v", os.Args[2])
		}

		localfs := flag.NewFlagSet(LOCAL, flag.ExitOnError)
		localhost := localfs.String("localhost", "localhost", "the local host to export")
		to := localfs.String("to", "me", "address of remote server")
		remote_port := localfs.Int("port", 8080, "port of remote server")

		if err := localfs.Parse(subArgs); err != nil {
			log.Fatal(err)
		}
		gopher.NewClient(*localhost, *to, localPort, *remote_port)
		// msg := gopher.HelloMessage {}
		// client.Send(msg)
	case SERVER:
		server := gopher.NewServer()
		server.Listen()
	default:
		log.Fatal("unrecognized command")
	}
}
