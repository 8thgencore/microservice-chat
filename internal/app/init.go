package app

import (
	"context"
	"log"

	"github.com/8thgencore/microservice-chat/internal/app/provider"
	"github.com/8thgencore/microservice-chat/internal/app/security"
	"github.com/8thgencore/microservice-chat/internal/config"
	"github.com/8thgencore/microservice-chat/internal/interceptor"
	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
	"github.com/8thgencore/microservice-common/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
	// Setup credentials
	creds, err := security.LoadServerCredentials(a.cfg.TLS.CertPath, a.cfg.TLS.KeyPath)
	if err != nil {
		log.Fatal(err)
	}

	c := a.serviceProvider.InterceptorClient()
	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			interceptor.LogInterceptor,
			interceptor.ValidateInterceptor,
			c.PolicyInterceptor,
		),
	)

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(a.grpcServer)

	// Register service with corresponded interface.
	chatv1.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImpl(ctx))

	return nil
}
