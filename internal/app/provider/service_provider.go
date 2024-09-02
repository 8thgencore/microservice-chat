package provider

import (
	"context"

	"github.com/8thgencore/microservice-chat/internal/config"
	"github.com/8thgencore/microservice-chat/internal/repository"
	"github.com/8thgencore/microservice-chat/pkg/db"

	chatRepository "github.com/8thgencore/microservice-chat/internal/repository/chat"
	logRepository "github.com/8thgencore/microservice-chat/internal/repository/log"
)

type ServiceProvider struct {
	Config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
}

func NewServiceProvider(config *config.Config) *ServiceProvider {
	return &ServiceProvider{
		Config: config,
	}
}

// Repository
func (s *ServiceProvider) UserRepository(ctx context.Context) repository.ChatRepository {
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
