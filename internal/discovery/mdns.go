package discovery

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	// ServiceName is the mDNS service identifier. 
	// It usually follows the format "_service._protocol"
	ServiceName = "_filetransfer._tcp"
	Domain      = "local."
)

// RegisterService starts broadcasting this application's presence on the network.
func RegisterService(instanceName string, port int) (*zeroconf.Server, error) {
	// Register a service by providing a human-readable instance name, 
	// the service name, the domain, and the port it is running on.
	// We can also pass optional text records (key-value pairs) for extra metadata.
	server, err := zeroconf.Register(instanceName, ServiceName, Domain, port, []string{"txtv=1", "app=myfiletransfer"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to register mDNS service: %w", err)
	}

	log.Printf("mDNS Service registered: %s (%s) on port %d\n", instanceName, ServiceName, port)
	return server, nil
}

// DiscoverPeers searches for other devices broadcasting our service on the network.
func DiscoverPeers(timeout time.Duration) ([]*zeroconf.ServiceEntry, error) {
	// Resolver to find services
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resolver: %w", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	var foundPeers []*zeroconf.ServiceEntry

	// Collect discovered peers from the channel
	go func() {
		for entry := range entries {
			log.Printf("Found peer: %s at %s:%d\n", entry.Instance, entry.AddrIPv4[0], entry.Port)
			foundPeers = append(foundPeers, entry)
		}
	}()

	// Context with timeout to stop browsing after a while
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Start browsing for the service
	log.Printf("Browsing for mDNS services (%s) for %v...\n", ServiceName, timeout)
	err = resolver.Browse(ctx, ServiceName, Domain, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to browse: %w", err)
	}

	// Wait for the timeout context to expire (meaning browse is done)
	<-ctx.Done()

	return foundPeers, nil
}
