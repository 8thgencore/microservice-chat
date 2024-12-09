package provider

import (
	"context"
	"log"

	"github.com/8thgencore/microservice-chat/internal/app/security"
	"github.com/8thgencore/microservice-chat/internal/client/rpc"
	"github.com/8thgencore/microservice-chat/internal/config"
	"github.com/8thgencore/microservice-chat/internal/delivery/chat"
	"github.com/8thgencore/microservice-chat/internal/interceptor"
	"github.com/8thgencore/microservice-chat/internal/repository"
	"github.com/8thgencore/microservice-chat/internal/service"
	"github.com/8thgencore/microservice-common/pkg/db"
	"google.golang.org/grpc"

	accessv1 "github.com/8thgencore/microservice-auth/pkg/pb/access/v1"
	rpcAuth "github.com/8thgencore/microservice-chat/internal/client/rpc/auth"

	chatRepository "github.com/8thgencore/microservice-chat/internal/repository/chat"
	logRepository "github.com/8thgencore/microservice-chat/internal/repository/log"
	messagesRepository "github.com/8thgencore/microservice-chat/internal/repository/messages"
	chatService "github.com/8thgencore/microservice-chat/internal/service/chat"
)

// ServiceProvider is a struct that provides access to various services and repositories.
type ServiceProvider struct {
	Config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	authClient rpc.AuthClient

	interceptorClient *interceptor.Client

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

// AuthClient creates a new instance of AuthClient
func (s *ServiceProvider) AuthClient() rpc.AuthClient {
	cfg := s.Config.AuthClient

	// Return existing client if already initialized
	if s.authClient != nil {
		return s.authClient
	}

	// Setup credentials
	creds, err := security.LoadClientCredentials(cfg.CertPath)
	if err != nil {
		log.Fatal(err)
	}

	// Establish gRPC connection with context
	conn, err := grpc.NewClient(
		cfg.Address(),
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("failed to connect to authentication service: %v", err)
	}

	// Initialize the auth client
	s.authClient = rpcAuth.NewAuthClient(accessv1.NewAccessV1Client(conn))

	return s.authClient
}

// InterceptorClient returns an instance of interceptor.Client.
func (s *ServiceProvider) InterceptorClient() *interceptor.Client {
	if s.interceptorClient == nil {
		s.interceptorClient = &interceptor.Client{
			Client: s.AuthClient(),
		}
	}

	return s.interceptorClient
}

// ChatRepository returns a chat repository.
func (s *ServiceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.chatRepository
}

// MessagesRepository returns a message repository.
func (s *ServiceProvider) MessagesRepository(ctx context.Context) repository.MessagesRepository {
	if s.messagesRepository == nil {
		s.messagesRepository = messagesRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.messagesRepository
}

// LogRepository returns a log repository.
func (s *ServiceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DatabaseClient(ctx))
	}
	return s.logRepository
}

// ChatService returns a chat service.
func (s *ServiceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.MessagesRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
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
