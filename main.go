// main.go
package main

import (
	"context"
	"log"
	"net"
	"time"
	"udp_iaasd/etcd"
	"udp_iaasd/libvirtctl"
	pb "udp_iaasd/proto"

	"google.golang.org/grpc"
)

const version = "0.0.1"

type server struct {
    pb.UnimplementedUtilsServer
    etcd *etcd.EtcdClient
}

// Ping
func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PongResponse, error) {
    log.Printf("Received: %v", in.GetMessage())
    return &pb.PongResponse{
        Message: "pong",
        Timestamp: time.Now().Unix(),
    }, nil
}

func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
    s.etcd.Put(ctx, "version", "super")
    res, err := s.etcd.Get(ctx, "version")
    if err != nil {
        log.Printf("Error getting value: %v", err)
    }
    print(res)
    return &pb.VersionResponse{
        Version: version,
    }, nil
}

func main() {
    // 1. initailize etcd client
    etcdClient := etcd.GetClient()
    defer etcdClient.Close()

    // 2. initailize libvirt client
    config := libvirtctl.DefaultConfig() // no reconnect
    conn := libvirtctl.GetInstance(config)

    if err := conn.Connect(); err != nil {
        log.Fatalf("failed to connect to libvirt: %v", err)
    }
    defer conn.Close()

    // Check if connected
    if conn.IsConnected() {
        log.Println("Connected to libvirt daemon")
    } else {
        log.Println("Not connected to libvirt daemon")
    }

    // // test etcd client
    // ctx := context.Background()
    // etcdClient.Put(ctx, "key", "value")
    // resp, err := etcdClient.Get(ctx, "key")
    // print(resp)
    // if err != nil {
    //     log.Printf("Error getting value: %v", err)
    // }
    // etcdClient.Delete(ctx, "key")



    // 2. start grpc server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterUtilsServer(s, &server{
        etcd: etcdClient,
    })

    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}