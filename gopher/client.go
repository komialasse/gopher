package gopher

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Message interface {
}

type HelloMessage struct {
	ForwardPort int
}

type AcceptMessage struct {
	uuid int
}

type Client struct {
	enc *gob.Encoder
	to string
	localhost string
	local_port int
	remote_port int
}

func (c *Client) Send(msg Message) {
	err := c.enc.Encode(&msg)
	if err != nil {
		panic(err)
	}
}

func NewClient(localhost, to string, local_port, remote_port int) *Client {
		addr := fmt.Sprintf("%s:%d", to, remote_port)
		log.Printf("client connecting to addr = %v\n", addr)
		conn, err := net.Dial("tcp", addr)
		enc := gob.NewEncoder(conn)
		if err != nil {
			panic(err)
		}
		return &Client{
			enc,
			to,
			localhost,
			local_port,
			remote_port,
		}
}