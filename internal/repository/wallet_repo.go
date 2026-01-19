package repo

import (
	"context"
	"learn-gin/internal/model"
)

type WalletRepo interface {
	Create(ctx context.Context, db DBTX, wallet *model.Wallet) error
}

type walletRepo struct {}

func NewWalletRepo() WalletRepo {
	return &walletRepo{}
}

func (r *walletRepo) Create(ctx context.Context, db DBTX, wallet *model.Wallet) error {
	query := `insert into wallets (user_id, balance, created_at, updated_at) values(?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt)
	return err
}