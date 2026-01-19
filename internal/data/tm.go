package data

import (
	"context"
	"database/sql"
	repo "learn-gin/internal/repository"
)

type TransactionManager interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context, tx repo.DBTX) error) error
	GetDB() repo.DBTX
}

type tm struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &tm{db: db}
}

func (t *tm) GetDB() repo.DBTX {
	return t.db
}

func (t *tm) ExecTx(ctx context.Context, fn func(ctx context.Context, tx repo.DBTX) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}