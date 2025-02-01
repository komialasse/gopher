package gopher

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	conn       *net.Conn
	stream     *Stream
	to         string
	localHost  string
	localPort  int
	remotePort int
}

func (c *Client) Send(msg Message) {
	err := c.stream.enc.Encode(&msg)
	if err != nil {
		panic(err)
	}
}

func NewStream(conn net.Conn) *Stream {
	return &Stream{&conn,
		gob.NewEncoder(conn),
		gob.NewDecoder(conn),
	}
}

func copy(closer chan struct{}, dst io.Writer, src io.ReadCloser) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{} // connection is closed, send signal to stop proxy
}

func proxy(local, remote *net.Conn) {
	closer := make(chan struct{}, 2)
	go copy(closer, *local, *remote)
	go copy(closer, *remote, *local)
	<-closer
}

func (c *Client) handleConn(Id uuid.UUID) {
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.to, DEFAULT_PORT))
	if err != nil {
		panic(err)
	}
	stream := NewStream(remoteConn)
	var accept Message = Message(Accept{Id})
	stream.enc.Encode(&accept)

	localConn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.localHost, c.localPort))
	if err != nil {
		panic(err)
	}

	go proxy(&localConn, &remoteConn)
}

func (c *Client) Listen() {
	defer (*c.conn).Close()

	for {
		var m Message
		err := c.stream.dec.Decode(&m)
		if err != nil {
			panic(err)
		}
		switch msg := m.(type) {
		case Connect:
			go c.handleConn(msg.Id)
		}
	}
}

func (c *Client) RemotePort() int {
	return c.remotePort
}

func NewClient(localHost, to string, localPort, port int) *Client {
	addr := fmt.Sprintf("%s:%d", to, DEFAULT_PORT)
	log.Printf("client connecting to %v\n", addr)
	conn, err := net.Dial("tcp", addr)
	stream := NewStream(conn)
	if err != nil {
		panic(err)
	}

	handshake := func() int {
		var hello Message = Hello{ Port: port }
		stream.enc.Encode(&hello)
		var m Message
		stream.dec.Decode(&m)
		switch msg := m.(type) {
		case Hello:
			return msg.Port
		default:
			panic("remote server did not respond with forward port")
		}
	}
	remotePort := handshake()

	return &Client{
		&conn,
		stream,
		to,
		localHost,
		localPort,
		remotePort,
	}
}
