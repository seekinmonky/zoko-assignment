package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/seekinmonky/zoko-messaging/config"
	"github.com/seekinmonky/zoko-messaging/db"
	"github.com/seekinmonky/zoko-messaging/handlers"
	"github.com/seekinmonky/zoko-messaging/queue"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database, err := db.NewDB(cfg.GetPostgresConnStr())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize tables
	if err := database.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	// Initialize Redis queue
	redisQueue := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPassword, database)

	// Start the worker to process messages from the queue
	redisQueue.StartWorker()

	// Initialize handlers
	messageHandler := handlers.NewMessageHandler(database, redisQueue)

	// Initialize router
	router := gin.Default()

	// Define routes
	router.POST("/messages", messageHandler.SendMessage)
	router.GET("/messages", messageHandler.GetConversation)
	router.PATCH("/messages/:message_id/read", messageHandler.MarkMessageAsRead)

	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
