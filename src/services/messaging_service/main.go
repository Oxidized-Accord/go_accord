package main

import (
	"log"
	"net"

	pb "messaging_service/messaging/proto"
	processes "messaging_service/processes"

	db "GOGOGO/src/libs/db"

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
	srv := processes.NewMessagingServer(database)

	pb.RegisterMessagingServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
