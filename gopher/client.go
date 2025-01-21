package gopher 

type Message interface {
}

type Hello struct {
	forward_port int
}
