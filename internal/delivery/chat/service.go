package chat

import (
	"github.com/8thgencore/microservice-chat/internal/service"
	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
)

// ChatImplementation structure describes API layer.
type ChatImplementation struct {
	chatv1.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewChatImplementation creates new object of API layer.
func NewChatImplementation(chatService service.ChatService) *ChatImplementation {
	return &ChatImplementation{
		chatService: chatService,
	}
}
