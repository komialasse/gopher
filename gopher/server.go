package gopher

import (
	"io"
	"log"
	"net"
)

func handle(conn net.Conn) {
	io.Copy(conn, conn)

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
