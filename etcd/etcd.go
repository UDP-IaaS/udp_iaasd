package etcd

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
    client  *clientv3.Client
    mu      sync.RWMutex
    endpoints []string
}

var (
    etcdClient *EtcdClient
    once sync.Once
)

func GetClient() *EtcdClient {
    once.Do(func() {
        etcdClient = &EtcdClient{
            endpoints: []string{"localhost:2379"},
        }
        if err := etcdClient.Connect(); err != nil {
            log.Fatalf("Failed to connect to etcd: %v", err)
        }
    })
    return etcdClient
}

func (e *EtcdClient) Connect() error {
    e.mu.Lock()
    defer e.mu.Unlock()
	
	// remove old client if exists
	if e.client != nil {
        e.client.Close()
    }

    config := clientv3.Config{
        Endpoints:   e.endpoints,
        DialTimeout: 5 * time.Second,
    }

    client, err := clientv3.New(config)
    if err != nil {
        return err
    }

    e.client = client
    return nil
}

func (e *EtcdClient) Close() {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    if e.client != nil {
        e.client.Close()
    }
}

func (e *EtcdClient) UpdateEndpoints(endpoints []string) error {
	// input validation
	if len(endpoints) == 0 {
        return errors.New("endpoints cannot be empty")
    }
	
	config := clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    }

    newClient, err := clientv3.New(config)
    if err != nil {
        return err
    }


    e.mu.Lock()
    defer e.mu.Unlock()
    
	if e.client != nil {
        e.client.Close()
    }

	// update the client and endpoints only if the new client is successfully created
	e.client = newClient
    e.endpoints = endpoints
    return e.Connect()
}

func (e *EtcdClient) RemoveEndpoint(endpoint string) error {
    e.mu.Lock()
    defer e.mu.Unlock()

	// input validation
	if len(e.endpoints) <= 1 {
        return errors.New("cannot remove last endpoint")
    }

	//allocate memory for the new endpoints, since it's predictable
    updated := make([]string, 0, len(e.endpoints)-1)
    for _, ep := range e.endpoints {
        if ep != endpoint {
            updated = append(updated, ep)
        }
    }

	config := clientv3.Config{
        Endpoints:   updated,
        DialTimeout: 5 * time.Second,
    }

	//create a new client with the updated endpoints
    newClient, err := clientv3.New(config)
    if err != nil {
        return err
    }

	// end connecton to the old client
    if e.client != nil {
        e.client.Close()
    }

    e.client = newClient
    e.endpoints = updated
    return nil
}

// Use RLock unless you need to modify something
func (e *EtcdClient) Put(ctx context.Context, key, value string) error {
    e.mu.RLock()
    defer e.mu.RUnlock()

	if e.client == nil {
		return errors.New("client not connected")
	}

    _, err := e.client.Put(ctx, key, value)
    return err
}

func (e *EtcdClient) Get(ctx context.Context, key string) (string, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

	if e.client == nil {
		return "error: client not connected", errors.New("client not connected")
	}

    resp, err := e.client.Get(ctx, key)
    if err != nil {
        return "", err
    }

    if len(resp.Kvs) == 0 {
        return "", nil
    }

    return string(resp.Kvs[0].Value), nil
}

func (e *EtcdClient) Delete(ctx context.Context, key string) error {
    e.mu.RLock()
    defer e.mu.RUnlock()

    if e.client == nil {
        return errors.New("client not connected")
    }

    _, err := e.client.Delete(ctx, key)
    return err
}