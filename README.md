# Messaging Platform API

A simple backend API for text-based messaging with asynchronous message processing using a Redis queue.

## Features

- Send text messages asynchronously using a Redis queue
- Retrieve conversation history between two users
- Mark messages as read

## Requirements

- Go 1.16+
- PostgreSQL
- Redis

## Setup

1. Make sure PostgreSQL and Redis are running on your machine
2. Create a PostgreSQL database named `messages`
3. Set environment variables (optional, defaults are provided):
   ```
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=messages
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
   SERVER_PORT=8080
   ```

## Running the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/api/main.go
```

## API Endpoints

### Send a Message (Asynchronous with Queue)
```
POST /messages
```
Request Body:
```json
{
  "sender_id": "user123",
  "receiver_id": "user456",
  "content": "Hello, how are you?"
}
```

### Retrieve Conversation History
```
GET /messages?user1=user123&user2=user456
```

### Mark a Message as Read
```
PATCH /messages/{message_id}/read
```

## Implementation Details

- Uses Go with Gin framework for routing
- PostgreSQL for message storage
- Redis for message queue implementation
- Asynchronous message processing with a worker 