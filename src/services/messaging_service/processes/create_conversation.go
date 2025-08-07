package processes

import (
	"context"
	"time"

	pb "messaging_service/messaging/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateConversation(ctx context.Context, req *pb.Conversation) (*pb.Conversation, error) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	query := "INSERT INTO conversations (id, timestamp) VALUES ($1, $2)"
	_, err := s.Db.ExecContext(ctx, query, id, timestamp)
	if err != nil {
		return nil, status.Errorf(500, "failed to create conversation, %v", err)
	}
	return &pb.Conversation{Id: id, Timestamp: timestamp}, nil
}
