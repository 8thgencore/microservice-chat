package provider

import (
	"context"

	"github.com/8thgencore/microservice-common/pkg/closer"
	"github.com/8thgencore/microservice-common/pkg/db"
	"github.com/8thgencore/microservice-common/pkg/db/pg"
	"github.com/8thgencore/microservice-common/pkg/db/transaction"
	"github.com/8thgencore/microservice-common/pkg/logger"
	"go.uber.org/zap"
)

// DatabaseClient returns a database client.
// If the client has not been created yet, it creates a new one using the DSN from the configuration.
// It also checks if the database is reachable by pinging it.
// The client is closed when the application shuts down.
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

// TxManager returns a transaction manager.
// If the transaction manager has not been created yet, it creates a new one using the database client.
func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DatabaseClient(ctx).DB())
	}
	return s.txManager
}
