package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/seekinmonky/zoko-messaging/db"
	"github.com/seekinmonky/zoko-messaging/models"
)

const (
	messageQueue = "message_queue"
)

// RedisQueue represents a Redis-based message queue
type RedisQueue struct {
	client *redis.Client
	db     *db.DB
	ctx    context.Context
}

// NewRedisQueue creates a new Redis queue
func NewRedisQueue(addr string, password string, db *db.DB) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &RedisQueue{
		client: client,
		db:     db,
		ctx:    context.Background(),
	}
}

// EnqueueMessage adds a message to the queue
func (q *RedisQueue) EnqueueMessage(msg models.Message) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return q.client.RPush(q.ctx, messageQueue, msgJSON).Err()
}

// StartWorker starts a worker to process messages from the queue
func (q *RedisQueue) StartWorker() {
	go func() {
		for {
			// Block until a message is available
			result, err := q.client.BLPop(q.ctx, 0*time.Second, messageQueue).Result()
			if err != nil {
				log.Printf("Error popping message from queue: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// The result is a slice where the first element is the queue name and the second is the value
			if len(result) < 2 {
				continue
			}

			msgJSON := result[1]
			var msg models.Message
			if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Save the message to the database
			if err := q.db.SaveMessage(msg); err != nil {
				log.Printf("Error saving message to database: %v", err)
				// Requeue the message
				if requeueErr := q.EnqueueMessage(msg); requeueErr != nil {
					log.Printf("Failed to requeue message: %v", requeueErr)
				}
			}
		}
	}()
}
