package soc// {{{

import (
	"adolf-codler/file_transfer/internal/transfer"
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	DEFAULT_PORT="4242"
)// }}}

// startServer sets up the TCP listener and accepts connections
func StartServer(path []string, onConnect func())(error){
	IP := getIP()
	if IP==""{
		return fmt.Errorf("Error getting IP")
	}
	log.Printf("Waiting for reciever on %v ...", IP)
	listener, err := net.Listen("tcp", net.JoinHostPort("",DEFAULT_PORT))
	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	defer listener.Close()
	log.Printf("Listening to port %s",DEFAULT_PORT)

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
	log.Printf("Sending welcome message")
	num := strconv.Itoa(no)
	msg:=fmt.Sprintf("Hello from the multi-file Go TCP server! number %s\n",num)
	conn.Write([]byte(msg))
	log.Printf("Message sent")
	log.Printf("Sending file")
	transfer.SendFile(conn, path)
}

// startClient connects to the server and reads the message stream
func StartClient(path []string, ip string) {
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

func getIP() string{// {{{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}// }}}

