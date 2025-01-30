package libvirtctl

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"libvirt.org/go/libvirt"
)

const (
	defaultURI = "qemu:///system"
	defaultReconnectInterval = 5 * time.Second
)
type ConnectionConfig struct {
	EnableAutoReconnect bool
	ReconnectInterval   time.Duration
}

func DefaultConfig() *ConnectionConfig {
	return &ConnectionConfig{
		EnableAutoReconnect: false,
		ReconnectInterval:   defaultReconnectInterval,
	}
}

var (
	instance *Connection
	once     sync.Once
)

// Connection represents a connection to the libvirt daemon
type Connection struct {
	conn    *libvirt.Connect
	config  *ConnectionConfig
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// GetInstance returns the singleton instance of Connection
func GetInstance(config *ConnectionConfig) *Connection {
	once.Do(func() {
		if config == nil {
			config = DefaultConfig()
		}
		ctx, cancel := context.WithCancel(context.Background())
		instance = &Connection{
			config: config,
			ctx:    ctx,
			cancel: cancel,
		}
	})
	return instance
}

// Connect establishes a connection to the libvirt daemon
func (c *Connection) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil // Already connected
	}

	conn, err := libvirt.NewConnect(defaultURI)
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt: %w", err)
	}

	c.conn = conn

	if c.config.EnableAutoReconnect {
		go c.monitorConnection()
	}

	return nil
}

// Close safely closes the connection
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}

	// Cancel monitoring goroutine before closing
	c.cancel()

	if _, err := c.conn.Close(); err != nil {
		c.conn = nil
		return fmt.Errorf("failed to close connection: %w", err)
	}

	c.conn = nil
	return nil
}

// GetConnection returns the underlying libvirt connection
func (c *Connection) GetConnection() (*libvirt.Connect, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return nil, fmt.Errorf("not connected to libvirt daemon")
	}

	alive, err := c.conn.IsAlive()
	if err != nil || !alive {
		return nil, fmt.Errorf("connection is not alive")
	}

	return c.conn, nil
}

// IsConnected checks if the connection is active
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return false
	}

	alive, err := c.conn.IsAlive()
	return err == nil && alive
}

// monitorConnection handles connection monitoring and automatic reconnection
func (c *Connection) monitorConnection() {
    ticker := time.NewTicker(c.config.ReconnectInterval)
    defer ticker.Stop()

    for {
        select {
        case <-c.ctx.Done():
            return
        case <-ticker.C:
            if !c.IsConnected() {
                c.mu.Lock()
                if c.conn != nil {
					if _, err := c.conn.Close(); err != nil {
						log.Printf("Failed to close existing connection: %v", err)
					}
                    c.conn = nil
                }
                // Attempt to reconnect
                conn, err := libvirt.NewConnect(defaultURI)
                if err != nil {
                    log.Printf("Failed to reconnect to libvirt: %v", err)
                    c.mu.Unlock()
                    continue
                }
                c.conn = conn
                log.Printf("Successfully reconnected to libvirt")
                c.mu.Unlock()
            }
        }
    }
}
