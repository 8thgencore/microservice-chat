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

// ServiceProvider is a struct that provides access to various services and repositories.
type ServiceProvider struct {
	Config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	chatRepository     repository.ChatRepository
	messagesRepository repository.MessagesRepository
	logRepository      repository.LogRepository

	chatService service.ChatService

	chatImpl *chat.Implementation
}

// NewServiceProvider creates a new instance of ServiceProvider with the given configuration.
func NewServiceProvider(config *config.Config) *ServiceProvider {
	return &ServiceProvider{
		Config: config,
	}
}

// ChatRepository returns a chat repository.
func (s *ServiceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.chatRepository
}

// LogRepository returns a log repository.
func (s *ServiceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.logRepository
}

// ChatService returns a chat service.
func (s *ServiceProvider) ChatService(_ context.Context) service.ChatService {
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

// ChatImpl returns a chat api implementation.
func (s *ServiceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}
	return s.chatImpl
}
