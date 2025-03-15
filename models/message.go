package models

import (
	"time"
)

// Message represents a message in the system
type Message struct {
	MessageID  string    `json:"message_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
	Read       bool      `json:"read"`
}

// MessageRequest represents the request body for sending a message
type MessageRequest struct {
	SenderID   string `json:"sender_id" binding:"required"`
	ReceiverID string `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// ReadResponse represents the response for marking a message as read
type ReadResponse struct {
	Status string `json:"status"`
}
