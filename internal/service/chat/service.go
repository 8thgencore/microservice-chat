package chat

import (
	"sync"

	"github.com/8thgencore/microservice-chat/internal/model"
	"github.com/8thgencore/microservice-chat/internal/repository"
	"github.com/8thgencore/microservice-chat/internal/service"
	"github.com/8thgencore/microservice-common/pkg/db"
)

const messagesBuffer = 100

type chatService struct {
	chatRepository     repository.ChatRepository
	messagesRepository repository.MessagesRepository
	logRepository      repository.LogRepository
	txManager          db.TxManager

	channels   map[string]chan *model.Message
	mxChannels sync.RWMutex

	chats  map[string]*chat
	mxChat sync.RWMutex
}

type chat struct {
	streams map[string]model.Stream
	m       sync.RWMutex
}

// NewService creates new object of service layer.
func NewService(
	chatRepository repository.ChatRepository,
	messagesRepository repository.MessagesRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.ChatService {
	return &chatService{
		chatRepository:     chatRepository,
		messagesRepository: messagesRepository,
		logRepository:      logRepository,
		txManager:          txManager,
		chats:              make(map[string]*chat),
		channels:           make(map[string]chan *model.Message),
	}
}
