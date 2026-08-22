package soc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

// startServer sets up the TCP listener and accepts connections
const (
	DEFAULT_PORT=":4242"
)
func StartServer() {
	listener, err := net.Listen("tcp", DEFAULT_PORT)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening to port %s",DEFAULT_PORT)

	no:=0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		no+=1
		go HandleClient(conn, no)
	}
}

// handleClient manages individual client connections concurrently
func HandleClient(conn net.Conn, no int) {
	defer conn.Close()
	num := strconv.Itoa(no)
	msg:=fmt.Sprintf("Hello from the multi-file Go TCP server! number %s\n",num)
	conn.Write([]byte(msg))
}

// startClient connects to the server and reads the message stream
func StartClient() {
	fmt.Println("Connecting to server at 127.0.0.1:8080...")
	conn, err := net.Dial("tcp", "127.0.0.1:4242")
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
