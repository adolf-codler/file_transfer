package soc

import (
	"adolf-codler/file_transfer/internal/data_transfer"
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

const (
	DEFAULT_PORT=":4242"
)

// startServer sets up the TCP listener and accepts connections
func StartServer(path string) {
	var wg sync.WaitGroup
	log.Printf("Server started")
	listener, err := net.Listen("tcp", DEFAULT_PORT)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening to port %s",DEFAULT_PORT)

	no:=0
	log.Printf("Waiting for client")
	conn, err := listener.Accept()
	log.Printf("Client Accepted")
	if err != nil {
		log.Printf("Connection error: %v", err)
	}
	log.Printf("Handling Client")
	wg.Add(1)
	go HandleClient(conn, no, &wg, path)
	wg.Wait()
}

// handleClient manages individual client connections concurrently
func HandleClient(conn net.Conn, no int, wg *sync.WaitGroup, path string) {
	defer wg.Done()
	log.Printf("Sending welcome message")
	num := strconv.Itoa(no)
	msg:=fmt.Sprintf("Hello from the multi-file Go TCP server! number %s\n",num)
	conn.Write([]byte(msg))
	log.Printf("Message sent")
	log.Printf("Sending file")
	transfer.SendFile(conn, path)
}

// startClient connects to the server and reads the message stream
func StartClient(path string) {
	log.Printf("Connecting to server")
	conn, err := net.Dial("tcp", "127.0.0.1:4242")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	log.Printf("Waiting for welcome")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Read error: %v", err)
	}

	fmt.Print("Server response: ", message)
	log.Printf("Receiving file")
	transfer.RecvFile(conn, path)
	log.Printf("Received")
}
