package gopher

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

const DEFAULT_PORT = 8080 

type Message interface {
}

type Hello struct {
	Port int
}

type Accept struct {
	Id uuid.UUID
}

type Client struct {
	stream *Stream
	to string
	localHost string
	localPort int
	remotePort int
}

func (c *Client) Send(msg Message) {
	err := c.stream.enc.Encode(&msg)
	if err != nil {
		panic(err)
	}
}

func (c *Client) Listen() {
	var m Message
	fmt.Println("decoding...")
	err := c.stream.dec.Decode(&m)
	fmt.Println("done decoding")
	if err != nil {
		panic(err)
	}
	switch msg := m.(type) {
	case Connect:
		fmt.Printf("server sent connect with id %v\n", msg.Id)
	}
}


func NewClient(localHost, to string, localPort, Port int) *Client {
		addr := fmt.Sprintf("%s:%d", to, DEFAULT_PORT)
		log.Printf("client connecting to addr = %v\n", addr)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		
		
		handshake := func() int {
			enc := gob.NewEncoder(conn)
			dec := gob.NewDecoder(conn)
			var hello Message = Hello { Port }
			enc.Encode(&hello)

			var m Message
			dec.Decode(&m)
			switch msg := m.(type) {
			case Hello:
				return msg.Port
			default:
				panic("remote server did not respond with forward port")
			}
		}
		remotePort := handshake()
		// Connect to server at remote proxy.
		addr = fmt.Sprintf("%s:%d", to, remotePort)
		log.Printf("client connecting to addr = %v\n", addr)
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		stream := Stream{gob.NewEncoder(conn), gob.NewDecoder(conn)}

		return &Client{
			&stream,
			to,
			localHost,
			localPort,
			remotePort,
		}
}