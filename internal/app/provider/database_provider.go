package provider

import (
	"context"

	"github.com/8thgencore/microservice-chat/pkg/closer"
	"github.com/8thgencore/microservice-chat/pkg/db"
	"github.com/8thgencore/microservice-chat/pkg/db/pg"
	"github.com/8thgencore/microservice-chat/pkg/db/transaction"
	"github.com/8thgencore/microservice-chat/pkg/logger"
	"go.uber.org/zap"
)

func (s *ServiceProvider) DatabaseClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		c, err := pg.New(ctx, s.Config.Database.DSN())
		if err != nil {
			logger.Fatal("failed to create db client: ", zap.Error(err))
		}

		err = c.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("failed to ping database: ", zap.Error(err))
		}

		closer.Add(c.Close)

		s.dbClient = c
	}

	return s.dbClient
}

func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DatabaseClient(ctx).DB())
	}
	return s.txManager
}
