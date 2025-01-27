package gopher

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func handle(dec *gob.Decoder) {
	var m Message
	err := dec.Decode(&m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("type of m: %T\n", (m))
}

type Server struct {}

func (s *Server) Listen() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	defer ln.Close()

	log.Printf("server listening on %v", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		log.Println("connection!");
		if err != nil {
			panic(err)
		}
		decoder := gob.NewDecoder(conn)
		go handle(decoder)
	}
}

func NewServer() *Server {
	return &Server {}
}
