package main

import (
	"log"
	"net"

	db "GOGOGO/src/libs/db"
	pb "messaging_service/messaging/proto"
	server "messaging_service/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	database := db.Connect()
	defer database.Close()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	srv := server.NewMessagingServer(database)

	pb.RegisterMessagingServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
