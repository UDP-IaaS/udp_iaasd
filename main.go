// main.go
package main

import (
	"context"
	"log"
	"net"
	"time"
	pb "udp_iaasd/proto"

	"google.golang.org/grpc"
)

const version = "0.0.1"

type server struct {
    pb.UnimplementedUtilsServer
}

// Ping implements ping.PingPongServer
func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PongResponse, error) {
    log.Printf("Received: %v", in.GetMessage())
    return &pb.PongResponse{
        Message: "pong",
        Timestamp: time.Now().Unix(),
    }, nil
}

func (s *server) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
    return &pb.VersionResponse{
        Version: version,
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterUtilsServer(s, &server{})
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}