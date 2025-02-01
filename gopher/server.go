package gopher

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/google/uuid"
)

type Server struct {
	conns map[uuid.UUID]*Stream
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
			minPort := 1024 // 2^10
			maxPort := 65536 // 2^16
			port = rand.Intn(maxPort-minPort) + 1 + minPort
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
		port := GetPort(ln.Addr())
		var hello Message = Hello{ port }
		stream.enc.Encode(&hello)

		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}

			id := uuid.New()
			server.conns[id] = NewStream(conn)
			var connect Message = Connect{ id }
			stream.enc.Encode(&connect)
		}
	case Accept:
		// forward connection
		otherStream, ok := server.conns[msg.Id]
		if !ok {
			log.Println("missing connection")
		}
		delete(server.conns, msg.Id)
		proxy(stream.conn, otherStream.conn)
	default:
		log.Println("unrecognized message")
	}
}

func (s *Server) Listen(ctx context.Context) error {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%v", DEFAULT_PORT))

	if err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	log.Printf("server listening on %v", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <- ctx.Done():
				return ctx.Err()
			default: 
				panic(err)
			}
		}
		log.Println("connection!")

		stream := NewStream(conn)
		go s.handle(stream)
	}
}

func NewServer() *Server {
	conns := make(map[uuid.UUID]*Stream)
	return &Server{conns}
}
