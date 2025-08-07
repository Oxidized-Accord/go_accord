package processes

import (
	"context"
	"database/sql"

	pb "messaging_service/messaging/proto"

	"google.golang.org/grpc/status"
)

func (s *Server) GetMessage(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	query := `SELECT id, conversation_id, sender, content, timestamp FROM messages WHERE id = $1`

	var msg pb.Message
	err := s.Db.QueryRowContext(ctx, query, req.Id).Scan(
		&msg.Id,
		&msg.ConversationId,
		&msg.Sender,
		&msg.Content,
		&msg.Timestamp,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(404, "message not found")
		}
		return nil, status.Errorf(500, "failed to get message: %v", err)
	}

	return &pb.GetResponse{Message: &msg}, nil
}
