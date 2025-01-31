package gopher

import (
	"encoding/gob"
	"fmt"
	"io"
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
	conn 		*net.Conn
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

func (c *Client) handleConnection(Id uuid.UUID) {
	remoteConnection, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.to, DEFAULT_PORT))
	if err != nil {
		panic(err)
	}
	stream := NewStream(remoteConnection)
	var accept Message = Message(Accept{Id})
	stream.enc.Encode(&accept)

	localConnection, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.localHost, c.localPort))
	if err != nil {
		panic(err)
	}

	go proxy(&localConnection, &remoteConnection)
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
			log.Printf("server connect with id: %v", msg.Id)
			go c.handleConnection(msg.Id)
		}
	}
}

func NewClient(localHost, to string, localPort, Port int) *Client {
	addr := fmt.Sprintf("%s:%d", to, DEFAULT_PORT)
	log.Printf("client connecting to addr = %v\n", addr)
	conn, err := net.Dial("tcp", addr)
	stream := Stream{&conn, gob.NewEncoder(conn), gob.NewDecoder(conn)}
	if err != nil {
		panic(err)
	}

	handshake := func() int {
		var hello Message = Hello{Port}
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
		&stream,
		to,
		localHost,
		localPort,
		remotePort,
	}
}
