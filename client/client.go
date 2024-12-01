package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"dilshat/faraway/pkg"
)

func main() {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	conn, err := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatalf("Failed to set read deadline: %v", err)
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read from server: %v", err)
	}
	challenge := string(buf[:n])
	fmt.Println("Challenge received:", challenge)

	ss := strings.Split(challenge, "|")
	if len(ss) != 2 {
		log.Fatalf("Malformed task received: %s", challenge)
	}

	task, hashPrefix := ss[0], ss[1]
	fmt.Printf("Task: %s, hashPrefix: %s\n", task, hashPrefix)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	nonce, err := pkg.FindNonce(ctx, task, hashPrefix)
	if err != nil {
		log.Fatalf("Failed to find nonce: %v", err)
	}
	fmt.Printf("Nonce found: %s\n", nonce)

	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Fatalf("Failed to set write deadline: %v", err)
	}

	_, err = conn.Write([]byte(nonce))
	if err != nil {
		log.Fatalf("Failed to send nonce to server: %v", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Fatalf("Failed to set read deadline: %v", err)
	}

	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read response from server: %v", err)
	}
	fmt.Println("Response from server:", string(buf[:n]))
}
