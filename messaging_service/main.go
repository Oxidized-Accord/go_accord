package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "messaging_service/messaging/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedMessagingServiceServer
	messages      sync.Map
	senderIndex   map[string][]string
	receiverIndex map[string][]string
	mu            sync.RWMutex
}

func (s *server) SendMessage(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	msg := &pb.Message{
		Id:        id,
		Sender:    req.Sender,
		Receiver:  req.Receiver,
		Content:   req.Content,
		Timestamp: timestamp,
	}

	s.messages.Store(id, msg)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.senderIndex[req.Sender] = append(s.senderIndex[req.Sender], id)

	s.receiverIndex[req.Receiver] = append(s.receiverIndex[req.Receiver], id)

	return &pb.SendResponse{Message: msg}, nil

}

func (s *server) GetMessage(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	val, ok := s.messages.Load(req.Id)
	if !ok {
		return nil, status.Errorf(404, "message not found")
	}
	return &pb.GetResponse{Message: val.(*pb.Message)}, nil
}

func (s *server) ListMessagesBySender(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.mu.RLock()
	ids := s.senderIndex[req.User]
	s.mu.RUnlock()

	var results []*pb.Message

	for _, id := range ids {
		if val, ok := s.messages.Load(id); ok {
			results = append(results, val.(*pb.Message))
		}
	}
	return &pb.ListResponse{Messages: results}, nil
}

func (s *server) ListMessagesByReceiver(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.mu.RLock()
	ids := s.receiverIndex[req.User]
	s.mu.RUnlock()

	var results []*pb.Message

	for _, id := range ids {
		if val, ok := s.messages.Load(id); ok {
			results = append(results, val.(*pb.Message))
		}
	}
	return &pb.ListResponse{Messages: results}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMessagingServiceServer(grpcServer, &server{
		senderIndex:   make(map[string][]string),
		receiverIndex: make(map[string][]string)})
	reflection.Register(grpcServer)

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
