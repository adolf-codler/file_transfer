package discovery

import (
	"context"
	"log"
	"net"
	"time"
)

const (
	// MulticastAddress is a private multicast group and port
	MulticastAddress = "239.255.255.250:9999"
	// DiscoveryMessage is the payload sent to announce presence
	DiscoveryMessage = "FILE_TRANSFER_DISCOVER"
)

// StartUDPBroadcaster periodically sends multicast discovery messages to the network.
func StartUDPBroadcaster(ctx context.Context, broadcastInterval time.Duration) error {
	addr, err := net.ResolveUDPAddr("udp4", MulticastAddress)
	if err != nil {
		return err
	}

	// Dial the multicast address
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	go func() {
		defer conn.Close()
		ticker := time.NewTicker(broadcastInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, err := conn.Write([]byte(DiscoveryMessage))
				if err != nil {
					log.Printf("UDP broadcast error: %v\n", err)
				}
			}
		}
	}()

	return nil
}

// StartUDPListener listens for multicast UDP discovery messages from other peers.
func StartUDPListener(ctx context.Context, onPeerDiscovered func(ip string)) error {
	addr, err := net.ResolveUDPAddr("udp4", MulticastAddress)
	if err != nil {
		return err
	}

	// Listen on the multicast group
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Use a short deadline so the goroutine can check for context cancellation
				_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				n, remoteAddr, err := conn.ReadFromUDP(buf)
				if err != nil {
					// Timeout or other read error, just continue loop
					continue
				}

				msg := string(buf[:n])
				if msg == DiscoveryMessage {
					if onPeerDiscovered != nil {
						onPeerDiscovered(remoteAddr.IP.String())
					}
				}
			}
		}
	}()

	return nil
}
