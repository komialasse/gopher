package gopher

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func handle(conn net.Conn) {
	var m Message
	var b []byte
	conn.Read(b)
	err := json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}

type Server struct {}

func (s *Server) Listen() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		log.Println("connection!");
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func NewServer() *Server {
	return &Server {}
}
