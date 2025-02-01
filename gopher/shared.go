package gopher

import (
	"encoding/gob"
	"net"
	"strconv"

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

type Connect struct {
	Id uuid.UUID
}

type Stream struct {
	conn *net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
}

func RegisterMessages() {
	gob.Register(Hello{})
	gob.Register(Connect{})
	gob.Register(Accept{})
}

func GetPort(addr net.Addr) int {
	_, p, err := net.SplitHostPort(addr.String())
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		panic(err)
	}
	return port
}