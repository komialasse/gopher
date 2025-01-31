package gopher

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	"github.com/google/uuid"
)

type Stream struct {
	conn *net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
}

type Connect struct {
	Id uuid.UUID
}

func getListner(port int, c chan net.Listener) {
	bind := func(port int) (net.Listener, error) {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		return ln, err
	}

	if port > 0 {
		// Port is not default 0 value.
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
			port = rand.Intn(max-min) + 1 + min
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

func (server *Server) handle(stream *Stream) {
	var m Message
	err := stream.dec.Decode(&m)
	if err != nil {
		panic(err)
	}
	switch msg := m.(type) {
	case Hello:
		ch := make(chan net.Listener)
		go getListner(msg.Port, ch)
		ln := <-ch
		_, p, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			panic(err)
		}
		Port, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		var hello Message = Hello{Port}
		log.Printf("sending hello at port %v\n", Port)
		stream.enc.Encode(&hello)

		for {
			log.Printf("waiting for accept on addr: %v\n", ln.Addr().String())
			conn, err := ln.Accept()
			log.Println("done waiting")
			enc, dec := gob.NewEncoder(conn), gob.NewDecoder(conn)
			if err != nil {
				panic(err)
			}

			Id := uuid.New()
			server.conns[Id] = Stream{&conn, enc, dec}
			var connect Message = Connect{Id}
			log.Println("server sending connect on conn")
			stream.enc.Encode(&connect)
			log.Println("done sending connect")
		}
	default:
	}
}

type Server struct {
	conns map[uuid.UUID]Stream
}

func (s *Server) Listen() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	defer ln.Close()

	log.Printf("server listening on %v", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		log.Println("connection!")
		if err != nil {
			panic(err)
		}
		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)
		server := NewServer()
		stream := Stream{&conn, enc, dec}
		go server.handle(&stream)
	}
}

func NewServer() *Server {
	conns := make(map[uuid.UUID]Stream)
	return &Server{conns}
}
