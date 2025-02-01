package gopher_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gopher"
)

func getClient() string {
	// TODO: add remaining setup code.
	client := gopher.NewClient("localhost", "localhost", 5050, 8080)
	remoteAddress := fmt.Sprintf("localhost:%v", client.RemotePort())
	return remoteAddress
}

func startServer(ctx context.Context) *gopher.Server {
	// TODO: add remaining setup.
	server := gopher.NewServer()
	go server.Listen(ctx)
	return server
}

func TestSetupServer(t *testing.T) {
	ctx, shutdown := context.WithCancel(context.Background())

	server := startServer(ctx)

	go func() {
		time.Sleep(1 * time.Second)
		shutdown()
	}()

	if server == nil {
		t.Error("server failed to initialize")
	}

	address := getClient()
	fmt.Println("address is ", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(conn, "hello world")
}
