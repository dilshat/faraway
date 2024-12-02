package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"dilshat/faraway/pkg"
)

var quotes = []string{
	"The only limit to our realization of tomorrow is our doubts of today.",
	"Do what you can, with what you have, where you are.",
	"The future belongs to those who believe in the beauty of their dreams.",
	"Don't watch the clock; do what it does. Keep going.",
	"Success usually comes to those who are too busy to be looking for it.",
}

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	hashPrefix := os.Getenv("HASH_PREFIX")
	if hashPrefix == "" {
		hashPrefix = "abc"
	}

	runServer(hashPrefix, port)

}

func runServer(hashPrefix, port string) {
	// Set timeout for listening connection
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("Server is running on :" + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(hashPrefix, conn)
	}
}

// handleConnection handles a single connection by sending a task and receiving a solution.
func handleConnection(hashPrefix string, conn net.Conn) {
	defer conn.Close()

	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("Failed to set read deadline: %v", err)
		return
	}

	log.Printf("New connection from %s", conn.RemoteAddr())

	task, err := generateTask()
	if err != nil {
		log.Printf("Failed to generate task: %v", err)
		return
	}

	_, err = fmt.Fprintf(conn, "%s|%s", task, hashPrefix)
	if err != nil {
		log.Printf("Failed to send task: %v", err)
		return
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read solution: %v", err)
		return
	}

	solution := string(buf[:n])
	if !pkg.VerifySolution(task, solution, hashPrefix) {
		fmt.Fprintln(conn, "Invalid solution. Goodbye!")
		log.Printf("Invalid solution from %s", conn.RemoteAddr())
		return
	}

	quote := getRandomQuote()
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("Failed to set write deadline: %v", err)
		return
	}

	fmt.Fprintf(conn, "Well done! Here's your quote: \"%s\"\n", quote)
	log.Printf("Quote sent to %s", conn.RemoteAddr())
}

// generateTask generates a random 8-byte task and returns it as a hexadecimal string.
func generateTask() (string, error) {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

// getRandomQuote returns a random quote from a predefined list of quotes.
func getRandomQuote() string {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(quotes))))
	return quotes[index.Int64()]
}
