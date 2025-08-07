package processes

import (
	"context"

	pb "messaging_service/messaging/proto"

	"google.golang.org/grpc/status"
)

func (s *Server) ListMessages(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	query := `SELECT id, conversation_id, sender, content, timestamp FROM messages WHERE conversation_id = $1 ORDER BY timestamp DESC`
	rows, err := s.Db.QueryContext(ctx, query, req.ConversationId)

	if err != nil {
		return nil, status.Errorf(500, "failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []*pb.Message
	for rows.Next() {
		var msg pb.Message
		err := rows.Scan(&msg.Id, &msg.ConversationId, &msg.Sender, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, status.Errorf(500, "failed to scan message: %v", err)
		}
		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(500, "error iterating messages: %v", err)
	}

	return &pb.ListResponse{Messages: messages}, nil
}
