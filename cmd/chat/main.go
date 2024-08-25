package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/8thgencore/microservice-chat/pkg/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type server struct {
	pb.UnimplementedChatV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func (s *server) Create(ctx context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	fmt.Printf("Create chat with users: %+v\n", req.GetUsernames())
	return &pb.CreateChatResponse{Id: 1}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteChatRequest) (*pb.Empty, error) {
	fmt.Printf("Delete chat: %d\n", req.GetId())
	return &pb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.Empty, error) {
	fmt.Printf("Send message from %s: %s\n", req.GetFrom(), req.GetText())
	return &pb.Empty{}, nil
}
