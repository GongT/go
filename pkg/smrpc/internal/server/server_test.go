package server

import (
	"net"
	"testing"
	"time"

	"github.com/gongt/go/internal/myenv"
)

func Test_Server_Start(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	s := New("test-server")
	err := s.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Simulate a client connection
	path := s.SocketPath()
	conn, err := net.Dial("unix", path)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Wait for a moment to allow the server to handle the connection
	time.Sleep(1 * time.Second)
}
