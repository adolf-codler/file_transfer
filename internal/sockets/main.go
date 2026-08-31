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
	DEFAULT_PORT="4242"
)

// startServer sets up the TCP listener and accepts connections
func StartServer(path string, onConnect func()) {
	var wg sync.WaitGroup
	IP := getIP()
	if IP==""{
		log.Fatal("Error getting IP")
	}
	log.Printf("Server started on %s", IP)
	listener, err := net.Listen("tcp", net.JoinHostPort("",DEFAULT_PORT))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening to port %s",DEFAULT_PORT)

	no:=0
	log.Printf("Waiting for client")
	conn, err := listener.Accept()
	if onConnect != nil {
		onConnect()
	}
	log.Printf("Client Accepted")
	if err != nil {
		log.Printf("Connection error: %v", err)
	}
	log.Printf("Handling Client")
	wg.Add(1)
	go handleClient(conn, no, &wg, path)
	wg.Wait()
}

// handleClient manages individual client connections concurrently
func handleClient(conn net.Conn, no int, wg *sync.WaitGroup, path string) {
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
func StartClient(path string, ip string) {
	log.Printf("Connecting to server")
	conn, err := net.Dial("tcp", net.JoinHostPort(ip,DEFAULT_PORT))
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

func getIP() string{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		// Check if the address is an IP net and not a loopback (like 127.0.0.1)
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

