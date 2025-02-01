package main

import (
	"context"
	"encoding/gob"
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

	gob.Register(gopher.Hello{})
	gob.Register(gopher.Connect{})
	gob.Register(gopher.Accept{})

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
		to := localfs.String("to", "localhost", "address of remote server")
		remotePort := localfs.Int("port", 8081, "port of remote server")

		if err := localfs.Parse(subArgs); err != nil {
			log.Fatal(err)
		}
		client := gopher.NewClient(*localhost, *to, localPort, *remotePort)
		client.Listen()
	case SERVER:
		server := gopher.NewServer()
		server.Listen(context.Background())
	default:
		log.Fatal("unrecognized command")
	}
}
