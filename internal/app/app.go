package app

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/8thgencore/microservice-chat/internal/app/provider"
	"github.com/8thgencore/microservice-chat/internal/config"
	"github.com/8thgencore/microservice-chat/internal/interceptor"
	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
	"github.com/8thgencore/microservice-common/pkg/closer"
	"github.com/8thgencore/microservice-common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App structure contains main application structures.
type App struct {
	cfg *config.Config

	serviceProvider *provider.ServiceProvider
	grpcServer      *grpc.Server
}

// NewApp creates new App object.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

// Run executes the application.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1) // gRPC servers

	go func() {
		defer wg.Done()

		if err := a.runGrpcServer(); err != nil {
			log.Fatal("failed to run gRPC server: ", error.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	a.cfg = cfg

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(string(a.cfg.Env))
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = provider.NewServiceProvider(a.cfg)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LogInterceptor,
			interceptor.ValidateInterceptor,
		),
	)

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(a.grpcServer)

	// Register service with corresponded interface.
	chatv1.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImpl(ctx))

	return nil
}

func (a *App) runGrpcServer() error {
	cfg := a.serviceProvider.Config.GRPC

	logger.Info("gRPC server running on ", zap.String("address", cfg.Address()))

	lis, err := net.Listen(cfg.Transport, cfg.Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
