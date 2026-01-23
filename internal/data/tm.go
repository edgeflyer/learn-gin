package data

import (
	"context"
	// repo "learn-gin/internal/repository"

	"gorm.io/gorm"
)

type TransactionManager interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error
	GetDB() *gorm.DB
}

type tm struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &tm{db: db}
}

func (t *tm) GetDB() *gorm.DB {
	return t.db
}

func (t *tm) ExecTx(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}