package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seekinmonky/zoko-messaging/db"
	"github.com/seekinmonky/zoko-messaging/models"
	"github.com/seekinmonky/zoko-messaging/queue"
)

// MessageHandler handles message-related requests
type MessageHandler struct {
	db    *db.DB
	queue *queue.RedisQueue
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(db *db.DB, queue *queue.RedisQueue) *MessageHandler {
	return &MessageHandler{
		db:    db,
		queue: queue,
	}
}

// SendMessage handles the request to send a message
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req models.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new message
	msg := db.CreateMessage(req)

	// Enqueue the message for asynchronous processing
	if err := h.queue.EnqueueMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue message"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message_id": msg.MessageID, "status": "queued"})
}

// GetConversation handles the request to get a conversation between two users
func (h *MessageHandler) GetConversation(c *gin.Context) {
	user1 := c.Query("user1")
	user2 := c.Query("user2")

	if user1 == "" || user2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both user1 and user2 parameters are required"})
		return
	}

	messages, err := h.db.GetConversation(user1, user2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve conversation"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// MarkMessageAsRead handles the request to mark a message as read
func (h *MessageHandler) MarkMessageAsRead(c *gin.Context) {
	messageID := c.Param("message_id")

	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID is required"})
		return
	}

	if err := h.db.MarkMessageAsRead(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ReadResponse{Status: "read"})
}
