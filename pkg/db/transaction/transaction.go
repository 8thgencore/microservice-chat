package transaction

import (
	"context"
	"fmt"

	"github.com/8thgencore/microservice-chat/pkg/db"
	"github.com/8thgencore/microservice-chat/pkg/db/pg"
	"github.com/jackc/pgx/v5"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager creates transaction manager which implements db.TxManager interface.
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// Check for nested transactions: if there is nested transaction, then there is no need to create a new transaction.
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// Start new transaction.
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	// Save transaction in context.
	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		// Recover from panic.
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		// If there is an error, rollback transaction.
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = fmt.Errorf("errRollback: %w, original error: %v", errRollback, err)
			}
			return
		}

		// If there is no error, commit transaction.
		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("errCommit: %w", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("failed to execute transaction: %w", err)
	}
	return err
}
