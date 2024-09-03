package provider

import (
	"context"

	"github.com/8thgencore/microservice-chat/internal/config"
	"github.com/8thgencore/microservice-chat/internal/delivery/chat"
	"github.com/8thgencore/microservice-chat/internal/repository"
	"github.com/8thgencore/microservice-chat/internal/service"
	"github.com/8thgencore/microservice-chat/pkg/db"

	chatRepository "github.com/8thgencore/microservice-chat/internal/repository/chat"
	logRepository "github.com/8thgencore/microservice-chat/internal/repository/log"
	chatService "github.com/8thgencore/microservice-chat/internal/service/chat"
)

type ServiceProvider struct {
	Config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	chatRepository     repository.ChatRepository
	messagesRepository repository.MessagesRepository
	logRepository      repository.LogRepository

	chatService service.ChatService

	chatImpl *chat.ChatImplementation
}

func NewServiceProvider(config *config.Config) *ServiceProvider {
	return &ServiceProvider{
		Config: config,
	}
}

// Repository
func (s *ServiceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.chatRepository
}

func (s *ServiceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.logRepository
}

// Service
func (s *ServiceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.chatRepository,
			s.messagesRepository,
			s.logRepository,
			s.txManager,
		)
	}
	return s.chatService
}

// Implementation
func (s *ServiceProvider) ChatImpl(ctx context.Context) *chat.ChatImplementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewChatImplementation(s.ChatService(ctx))
	}
	return s.chatImpl
}
