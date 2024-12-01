package main

import (
	"context"
	"dilshat/faraway/pkg"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	os.Setenv("SERVER_PORT", "8080")
	go main()

	// Create a client to connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Read the challenge from the server
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}
	challenge := string(buf[:n])

	// Ensure the challenge contains the expected format
	if !strings.Contains(challenge, "|abc") {
		t.Errorf("Challenge does not contain expected prefix: %v", challenge)
	}

	ss := strings.Split(challenge, "|")
	if len(ss) != 2 {
		t.Fatalf("Malformed task received: %s", challenge)
	}

	task, hashPrefix := ss[0], ss[1]

	// Send a mock solution to the server

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	nonce, err := pkg.FindNonce(ctx, task, hashPrefix)
	if err != nil {
		t.Fatalf("Failed to find nonce: %v", err)
	}

	_, err = conn.Write([]byte(nonce))
	if err != nil {
		t.Fatalf("Failed to send nonce to server: %v", err)
	}

	// Read the response from the server
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}
	response := string(buf[:n])

	// Check that the response contains the expected "Well done!" message
	if !strings.Contains(response, "Well done!") {
		t.Errorf("Unexpected response from server: %v", response)
	}
}
