# mDNS (Zeroconf) Discovery in Go

mDNS (Multicast DNS) and DNS-SD (DNS Service Discovery) form the backbone of "Zero Configuration Networking" (Zeroconf). This is the exact same protocol suite used by Apple's Bonjour, Google's Chromecast, and Linux's Avahi to let devices find each other on a local network without any IP addresses being manually typed.

## How it works

1. **Service Registration (Announcing):**
   When your application starts, it registers a "Service" on the network. A service is defined by:
   - **Instance Name:** A human-readable name (e.g., "Alice's Macbook").
   - **Service Type:** A string defining the protocol and transport (e.g., `_filetransfer._tcp`).
   - **Domain:** Typically `local.`
   - **Port:** The TCP or UDP port your application is actually listening on for connections.
   - **TXT Records:** Optional key-value pairs for metadata (e.g., `version=1.0`).

   Your application broadcasts this registration to a special Multicast IP address (`224.0.0.251` for IPv4).

2. **Service Discovery (Browsing):**
   When a user wants to connect to a peer, the application sends out a query for the specific service type (e.g., `_filetransfer._tcp`) to that same multicast IP.
   Any device on the network that has registered that service responds with its IP address, port, and metadata.

## Using `zeroconf` in Go

The popular library `github.com/grandcat/zeroconf` makes this very easy.

### Registering a Service
```go
import "github.com/grandcat/zeroconf"

// Register the service (starts a background goroutine announcing presence)
server, err := zeroconf.Register(
    "MyFileTransfer",      // Instance Name
    "_filetransfer._tcp",  // Service Type
    "local.",              // Domain
    8080,                  // The actual port for the file transfer
    []string{"version=1"}, // TXT Records
    nil,                   // Network interfaces (nil means all)
)
defer server.Shutdown()
```

### Discovering a Service
```go
import "github.com/grandcat/zeroconf"

resolver, err := zeroconf.NewResolver(nil)

entries := make(chan *zeroconf.ServiceEntry)
go func() {
    for entry := range entries {
        fmt.Println("Found service at IP:", entry.AddrIPv4[0], "Port:", entry.Port)
    }
}()

// Browse for the service for 5 seconds
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err = resolver.Browse(ctx, "_filetransfer._tcp", "local.", entries)
<-ctx.Done()
```

## Advantages over Custom UDP
- **Collision Resolution:** If two devices use the name "Alice's Macbook", mDNS automatically renames one to "Alice's Macbook (2)".
- **TTL and Caching:** Devices cache the responses. When a device disconnects properly, it broadcasts a "goodbye" packet.
- **Interoperability:** You can use command-line tools like `dns-sd -B _filetransfer._tcp` on macOS or `avahi-browse` on Linux to debug and see your services running without needing custom clients!
