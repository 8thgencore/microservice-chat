package chat

import (
	"github.com/8thgencore/microservice-chat/internal/service"
	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
)

// Implementation structure describes API layer.
type Implementation struct {
	chatv1.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation creates new object of API layer.
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
