package server

import (
	"database/sql"
	pb "messaging_service/messaging/proto"
)

type ServerState struct {
	pb.UnimplementedMessagingServiceServer
	Db *sql.DB
}

func NewMessagingServer(db *sql.DB) *ServerState {
	return &ServerState{Db: db}
}
