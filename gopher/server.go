package gopher

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
)


func getListner(port int, c chan net.Listener) {
	bind := func(port int) (net.Listener, error) {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		return ln, err
	}

	if port > 0 {
		// Port is not default 0 value.
		fmt.Println("port is greater than 0")
		ln, err := bind(port)
		if err != nil {
			panic(err)
		}
		c <- ln
	} else {
		for range 100 {
			// Find a  port in some sample of ports
			min := 10
			max := 50
			port = rand.Intn(max - min) + 1 + min
			ln, err := bind(port)

			if err != nil {
				continue
			} else {
				c <- ln
				return 
			}
		}

		panic("unable to find port")
	}

}

func handle(enc *gob.Encoder, dec *gob.Decoder) {
	var m Message
	err := dec.Decode(&m)
	if err != nil {
		panic(err)
	}
	switch msg := m.(type) {
	case HelloMessage:
		ch := make(chan net.Listener)
		go getListner(msg.ForwardPort, ch)
		ln := <- ch
		fmt.Println("received ln")
		_, p, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			panic(err)
		}
		port, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		var hello Message = HelloMessage { ForwardPort: port}
		enc.Encode(&hello)
	default:
	}
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
		enc := gob.NewEncoder(conn) 
		dec := gob.NewDecoder(conn)
		go handle(enc, dec)
	}
}

func NewServer() *Server {
	return &Server {}
}
