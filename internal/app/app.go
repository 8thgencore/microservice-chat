package app

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/8thgencore/microservice-chat/internal/app/provider"
	"github.com/8thgencore/microservice-chat/internal/config"
	"github.com/8thgencore/microservice-common/pkg/closer"
	"github.com/8thgencore/microservice-common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
