# etcd Client Package

## Usage Guide

### Initialization

```go
import "udp_iaasd/etcd"

func main() {
    client := etcd.GetClient() 
    // Singleton instance
}
```

### Basic Operations

#### Put a Key-Value Pair

```go
err := client.Put(context.Background(), "foo", "bar")
if err != nil {
    log.Fatal(err)
}
```

#### Get a Value

```go
value, err := client.Get(context.Background(), "foo")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Value:", value)
```

#### Delete a Key

```go
err := client.Delete(context.Background(), "foo")
if err != nil {
    log.Fatal(err)
}
```

### Endpoint Management

#### Update Endpoints

```go
newEndpoints := []string{"node1:2379", "node2:2379", "node3:2379"}
err := client.UpdateEndpoints(newEndpoints)
if err != nil {
    log.Fatal("Failed to update endpoints:", err)
}
```

#### Remove an Endpoint

```go
err := client.RemoveEndpoint("node3:2379")
if err != nil {
    log.Fatal("Failed to remove endpoint:", err)
}
```

### Cleanup

```go
defer client.Close()
```

## Implementation Details

### Structure

* **Singleton Pattern**: Ensures single instance via `sync.Once`
* **Concurrency Safety**: Uses `sync.RWMutex` for coordinated access
* **Connection Management**: Automatic reconnection with endpoint changes

### Key Components

#### EtcdClient Struct

```go
type EtcdClient struct {
    client     *clientv3.Client
    mu         sync.RWMutex
    endpoints  []string
}
```

## Maintenance Notes

### Design Rationale

1. **Singleton Pattern**:
   * Prevents multiple clients competing for connections
   * Centralizes endpoint configuration management
2. **Mutex Usage**:
   * Protects against concurrent map writes in etcd client
   * Ensures atomic client swaps during reconnections
3. **Endpoint Validation**:
   * Maintains minimum cluster availability requirements
   * Prevents accidental configuration errors

### Extension Points

1. **Custom Configuration**:

```go
func NewCustomClient(endpoints []string, timeout time.Duration) *EtcdClient
```

2. **Enhanced Logging**:
   * Add debug logs for connection state changes
3. **Retry Logic**:
   * Implement exponential backoff for connection attempts

### Best Practices

1. **Context Usage**:
   * Always pass contexts with timeouts for production operations

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

2. **Endpoint Management**:
   * Monitor cluster health before removing endpoints
   * Prefer horizontal scaling over endpoint removal