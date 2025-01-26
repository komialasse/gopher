package gopher

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Message interface {
}

type HelloMessage struct {
	forward_port int
}

type AcceptMessage struct {
	uuid int
}

type Client struct {
	conn net.Conn
	to string
	localhost string
	local_port int
	remote_port int
}

func (c *Client) Send(msg Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	c.conn.Write(b)
}

func NewClient(localhost, to string, local_port, remote_port int) *Client {
		addr := fmt.Sprintf("%s:%d", to, remote_port)
		log.Printf("client connecting to addr = %v\n", addr)
		// conn, err := net.Dial("tcp", addr)
		// if err != nil {
		// 	panic(err)
		// }
		return &Client{
			nil,
			to,
			localhost,
			local_port,
			remote_port,
		}
}