package processes

import (
	"context"
	"database/sql"
	"time"

	pb "messaging_service/messaging/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedMessagingServiceServer
	db *sql.DB
}

func NewMessagingServer(db *sql.DB) *server {
	return &server{db: db}
}

func (s *server) SendMessage(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Insert message into database
	query := `INSERT INTO messages (id, conversation_id, sender, content, timestamp) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query, id, req.ConversationId, req.Sender, req.Content, timestamp)
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

func (s *server) GetMessage(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	query := `SELECT id, conversation_id, sender, content, timestamp FROM messages WHERE id = $1`

	var msg pb.Message
	err := s.db.QueryRowContext(ctx, query, req.Id).Scan(
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

func (s *server) ListMessages(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	query := `SELECT id, conversation_id, sender, content, timestamp FROM messages WHERE conversation_id = $1 ORDER BY timestamp DESC`
	rows, err := s.db.QueryContext(ctx, query, req.ConversationId)

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

/*
	func (s *server) ListMessagesBySender(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
		query := `SELECT id, sender, Sender, content, timestamp FROM messages WHERE Sender = $1 ORDER BY timestamp DESC`

		rows, err := s.db.QueryContext(ctx, query, req.ConversationId)
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
*/
func (s *server) CreateConversation(ctx context.Context, req *pb.Conversation) (*pb.Conversation, error) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	query := "INSERT INTO conversations (id, timestamp) VALUES ($1, $2)"
	_, err := s.db.ExecContext(ctx, query, id, timestamp)
	if err != nil {
		return nil, status.Errorf(500, "failed to create conversation, %v", err)
	}
	return &pb.Conversation{Id: id, Timestamp: timestamp}, nil
}
