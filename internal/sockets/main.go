package soc 

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// startServer sets up the TCP listener and accepts connections
func StartServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go HandleClient(conn)
	}
}

// handleClient manages individual client connections concurrently
func HandleClient(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Hello from the multi-file Go TCP server!\n"))
}

// startClient connects to the server and reads the message stream
func StartClient() {
	fmt.Println("Connecting to server at 127.0.0.1:8080...")
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Read error: %v", err)
	}

	fmt.Print("Server response: ", message)
}
