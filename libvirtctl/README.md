# libvirtctl

## Usage

### Basic Connection

Here's a simple example of how to use the package:

```go
package main

import (
    "log"
    "udp_iaasd/libvirtctl"
)

func main() {
    // Get the default connection instance
    conn := libvirtctl.GetInstance(nil)
    
    // Connect to libvirt daemon
    if err := conn.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Use the connection
    libvirtConn, err := conn.GetConnection()
    if err != nil {
        log.Fatalf("Failed to get connection: %v", err)
    }
    
    // Use libvirtConn for libvirt operations
}
```

### With Auto-reconnect

To enable automatic reconnection:

```go
package main

import (
    "log"
    "time"
    "udp_iaasd/libvirtctl"
)

func main() {
    // Configure connection with auto-reconnect
    config := &libvirtctl.ConnectionConfig{
        EnableAutoReconnect: true,
        ReconnectInterval:   10 * time.Second,
    }
    
    // Get connection instance with config
    conn := libvirtctl.GetInstance(config)
    
    // Connect to libvirt daemon
    if err := conn.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Connection will automatically attempt to reconnect if lost
}
```

### Check Connection Status

```go
package main

import (
    "log"
    "udp_iaasd/libvirtctl"
)

func main() {
    conn := libvirtctl.GetInstance(nil)
    
    if err := conn.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Check if connected
    if conn.IsConnected() {
        log.Println("Connected to libvirt daemon")
    } else {
        log.Println("Not connected to libvirt daemon")
    }
}
```

## Configuration Options

The `ConnectionConfig` struct provides the following configuration options:

```go
type ConnectionConfig struct {
    EnableAutoReconnect bool          // Enable/disable auto-reconnection
    ReconnectInterval   time.Duration // Interval between reconnection attempts
}
```

You can use `DefaultConfig()` to get the default configuration:
- `EnableAutoReconnect`: false
- `ReconnectInterval`: 5 seconds

## Thread Safety

The package is designed to be thread-safe and can be safely used across multiple goroutines. All operations that access the connection are protected by appropriate mutex locks.

## Best Practices

1. Always defer `Close()` after establishing a connection
2. Check connection status before performing operations
3. Handle errors appropriately in your application
4. Consider enabling auto-reconnect for long-running applications