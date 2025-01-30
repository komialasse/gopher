package gopher

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

const DEFAULT_PORT = 8080 

type Message interface {
}

type Hello struct {
	Port int
}

type Client struct {
	enc *gob.Encoder
	to string
	localHost string
	localPort int
	remotePort int
}

func (c *Client) Send(msg Message) {
	err := c.enc.Encode(&msg)
	if err != nil {
		panic(err)
	}
}


func NewClient(localHost, to string, localPort, Port int) *Client {
		addr := fmt.Sprintf("%s:%d", to, DEFAULT_PORT)
		log.Printf("client connecting to addr = %v\n", addr)
		conn, err := net.Dial("tcp", addr)
		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)
		if err != nil {
			panic(err)
		}
		


		handshake := func() int {
			var hello Message = Hello { Port }
			enc.Encode(&hello)

			var m Message
			dec.Decode(&m)
			fmt.Printf("m: %v\n", m)
			switch msg := m.(type) {
			case Hello:
				return msg.Port
			default:
				panic("remote server did not respond with forward port")
			}
		}
		remotePort := handshake()

		return &Client{
			enc,
			to,
			localHost,
			localPort,
			remotePort,
		}
}