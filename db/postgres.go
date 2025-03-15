package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/seekinmonky/zoko-messaging/models"
)

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// InitTables creates the necessary tables if they don't exist
func (db *DB) InitTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		message_id VARCHAR(36) PRIMARY KEY,
		sender_id VARCHAR(255) NOT NULL,
		receiver_id VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		read BOOLEAN NOT NULL DEFAULT false
	);
	CREATE INDEX IF NOT EXISTS idx_messages_sender_receiver ON messages(sender_id, receiver_id);
	CREATE INDEX IF NOT EXISTS idx_messages_receiver_sender ON messages(receiver_id, sender_id);
	`

	_, err := db.Exec(query)
	return err
}

// SaveMessage saves a message to the database
func (db *DB) SaveMessage(msg models.Message) error {
	query := `
	INSERT INTO messages (message_id, sender_id, receiver_id, content, timestamp, read)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.Exec(query, msg.MessageID, msg.SenderID, msg.ReceiverID, msg.Content, msg.Timestamp, msg.Read)
	return err
}

// GetConversation retrieves the conversation between two users
func (db *DB) GetConversation(user1, user2 string) ([]models.Message, error) {
	query := `
	SELECT message_id, sender_id, receiver_id, content, timestamp, read
	FROM messages
	WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY timestamp ASC
	`

	rows, err := db.Query(query, user1, user2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp, &msg.Read); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// MarkMessageAsRead marks a message as read
func (db *DB) MarkMessageAsRead(messageID string) error {
	query := `
	UPDATE messages
	SET read = true
	WHERE message_id = $1
	`

	result, err := db.Exec(query, messageID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message with ID %s not found", messageID)
	}

	return nil
}

// CreateMessage creates a new message with a unique ID
func CreateMessage(req models.MessageRequest) models.Message {
	return models.Message{
		MessageID:  uuid.New().String(),
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		Timestamp:  time.Now().UTC(),
		Read:       false,
	}
}
