package cus_soc// {{{

import (
	"adolf-codler/file_transfer/internal/transfer"
	"fmt"
	"log"
	"net"
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
	log.Printf("Sending file")
	err=transfer.SendFile(conn, path)
	if err!=nil{
		return fmt.Errorf("%w", err)
	}
	return nil
}

// startClient connects to the server and reads the message stream
func StartClient(path []string, ip string)error {
	log.Printf("Connecting to server")
	conn, err := net.Dial("tcp", net.JoinHostPort(ip,DEFAULT_PORT))
	if err != nil {
		return fmt.Errorf("Dial Error: %w", err)
	}
	log.Printf("Receiving file")
	err=transfer.RecvFile(conn)
	if err!=nil{
		return fmt.Errorf("%w", err)
	}
	log.Printf("Received")
	return nil
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

