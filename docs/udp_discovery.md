# UDP Discovery in Go

UDP discovery is a common technique used to find services or peers on a local network without prior knowledge of their IP addresses. This is typically done using **UDP Broadcast** or **UDP Multicast**. 

Here is a guide on how to implement a simple UDP Broadcast discovery mechanism in Go.

## 1. The Broadcaster (Client/Sender)

The broadcaster sends a message to the broadcast IP address (e.g., `255.255.255.255` for the local network) on a specific port. Any device listening on that port can receive the message.

```go
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// The broadcast address and the specific port to use
	broadcastAddr := "255.255.255.255:9999"
	
	addr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	if err != nil {
		panic(err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Periodically send discovery packets
	for {
		msg := []byte("DISCOVER_PEER")
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("Error sending broadcast:", err)
		} else {
			fmt.Println("Sent discovery packet...")
		}
		time.Sleep(3 * time.Second) // Wait before sending again
	}
}
```

## 2. The Listener (Server/Receiver)

The listener binds to the specific UDP port and listens for incoming messages. When it receives a discovery message, it can extract the sender's IP address and respond directly to them, establishing a connection.

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	listenPort := ":9999"

	addr, err := net.ResolveUDPAddr("udp4", listenPort)
	if err != nil {
		panic(err)
	}

	// Listen for UDP packets on the specified port
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Listening for discovery packets on port 9999...")

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		
		msg := string(buf[:n])
		fmt.Printf("Received '%s' from %s\n", msg, remoteAddr)

		// If the message is a valid discovery request, we can respond
		if msg == "DISCOVER_PEER" {
			response := []byte("PEER_ACK")
			_, err = conn.WriteToUDP(response, remoteAddr)
			if err != nil {
				fmt.Println("Error sending response:", err)
			} else {
				fmt.Printf("Sent acknowledgment to %s\n", remoteAddr)
			}
		}
	}
}
```

## Key Considerations

1. **Firewalls**: System firewalls often block incoming UDP traffic by default. Make sure to allow traffic on your chosen discovery port.
2. **Multicast vs Broadcast**: 
   - **Broadcast** (`255.255.255.255`) is simpler but is often restricted by routers and only works on the immediate local subnet.
   - **Multicast** (e.g., `224.0.0.1`) allows a group of devices to listen to a specific IP address without flooding the entire network. Go supports this via `net.ListenMulticastUDP`.
3. **Payload Structure**: In a real-world scenario, instead of plain text, consider using a structured payload (like JSON or Protocol Buffers) that includes information such as the sender's ID, services offered, or a TCP port for a follow-up reliable connection.
