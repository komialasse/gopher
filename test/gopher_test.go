package gopher_test

import (
	"context"
	"fmt"
	"net"
	"slices"
	"sync"
	"testing"

	"time"

	"github.com/gopher"
)

const REMOTE_PORT = 8081

const helloWorld = "hello world"

func getClient(ctx context.Context) (net.Listener, string) {
	listener, err := net.Listen("tcp", ":0");
	if err != nil {
		panic(err)
	}
	port := gopher.GetPort(listener.Addr())
	client := gopher.NewClient("localhost", "localhost", port, REMOTE_PORT)
	remoteAddress := fmt.Sprintf("localhost:%v", client.RemotePort())
	go client.Listen(ctx)
	return listener, remoteAddress
}

func startServer(ctx context.Context) *gopher.Server {
	server := gopher.NewServer()
	go server.Listen(ctx)
	return server
}

func setup() {
	gopher.RegisterMessages()
}

func TestSetupServer(t *testing.T) {
	setup()
	
	var wg sync.WaitGroup
	ctx, shutdown := context.WithCancel(context.Background())

	server := startServer(ctx)

	go func() {
		time.Sleep(1 * time.Second)
		shutdown()
	}()

	if server == nil {
		t.Error("server failed to initialize")
	}

	listener, address := getClient(ctx)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 11)
		c.Read(buf)
		if !slices.Equal(buf, []byte(helloWorld)) {
			t.Errorf("buf = %v, wanted %v", string(buf), helloWorld)
		}
		wg.Done()
	}(&wg)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte(helloWorld))

	wg.Wait()
}
