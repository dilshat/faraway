package main

import (
	"net"
	"testing"
)

// Mock server to simulate the challenge-response server
func startMockServer(t *testing.T) net.Listener {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			defer conn.Close()

			// Send a mock challenge
			challenge := "task123|abc"
			if _, err = conn.Write([]byte(challenge)); err != nil {
				t.Logf("Failed to write task to client: %v", err)
			}

			// Read the client's response
			buf := make([]byte, 256)
			_, err = conn.Read(buf)
			if err != nil {
				t.Logf("Failed to read client response: %v", err)
			}

			// Send a success message
			if _, err = conn.Write([]byte("Well done!")); err != nil {
				t.Logf("Failed to write response to client: %v", err)
			}
		}
	}()

	return ln
}

func TestClient(t *testing.T) {
	server := startMockServer(t)
	defer server.Close()

	runClient("localhost", "8080")
}
