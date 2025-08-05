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
	query := `INSERT INTO messages (id, sender, receiver, content, timestamp) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query, id, req.Sender, req.Receiver, req.Content, timestamp)
	if err != nil {
		return nil, status.Errorf(500, "failed to save message: %v", err)
	}

	msg := &pb.Message{
		Id:        id,
		Sender:    req.Sender,
		Receiver:  req.Receiver,
		Content:   req.Content,
		Timestamp: timestamp,
	}

	return &pb.SendResponse{Message: msg}, nil
}

func (s *server) GetMessage(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	query := `SELECT id, sender, receiver, content, timestamp FROM messages WHERE id = $1`

	var msg pb.Message
	err := s.db.QueryRowContext(ctx, query, req.Id).Scan(
		&msg.Id,
		&msg.Sender,
		&msg.Receiver,
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

func (s *server) ListMessagesBySender(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	query := `SELECT id, sender, receiver, content, timestamp FROM messages WHERE sender = $1 ORDER BY timestamp DESC`

	rows, err := s.db.QueryContext(ctx, query, req.User)
	if err != nil {
		return nil, status.Errorf(500, "failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []*pb.Message
	for rows.Next() {
		var msg pb.Message
		err := rows.Scan(&msg.Id, &msg.Sender, &msg.Receiver, &msg.Content, &msg.Timestamp)
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

func (s *server) ListMessagesByReceiver(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	query := `SELECT id, sender, receiver, content, timestamp FROM messages WHERE receiver = $1 ORDER BY timestamp DESC`

	rows, err := s.db.QueryContext(ctx, query, req.User)
	if err != nil {
		return nil, status.Errorf(500, "failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []*pb.Message
	for rows.Next() {
		var msg pb.Message
		err := rows.Scan(&msg.Id, &msg.Sender, &msg.Receiver, &msg.Content, &msg.Timestamp)
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
