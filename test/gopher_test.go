package gopher_test

import (
	"context"
	"testing"
	"time"

	"github.com/gopher"
)

func getClient() *gopher.Client {
	// TODO: add remaining setup code.
	return gopher.NewClient("localhost", "localhost", 5050, 8080)
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

}
