package model

import (
	"time"

	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
)

// Chat type is the main structure for chat.
type Chat struct {
	ID        string
	Usernames []string
}

// Message type is the main structure for user message.
type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}

// Stream is the wrapper for gRPC stream interface.
type Stream interface {
	chatv1.ChatV1_ConnectServer
}
