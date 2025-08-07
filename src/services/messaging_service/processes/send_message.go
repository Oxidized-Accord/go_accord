package processes

import (
	"context"
	"time"

	pb "messaging_service/messaging/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func (s *Server) SendMessage(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Insert message into database
	query := `INSERT INTO messages (id, conversation_id, sender, content, timestamp) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.Db.ExecContext(ctx, query, id, req.ConversationId, req.Sender, req.Content, timestamp)
	if err != nil {
		return nil, status.Errorf(500, "failed to save message: %v", err)
	}

	msg := &pb.Message{
		Id:             id,
		ConversationId: req.ConversationId,
		Sender:         req.Sender,
		Content:        req.Content,
		Timestamp:      timestamp,
	}

	return &pb.SendResponse{Message: msg}, nil
}
