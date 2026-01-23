package repo

import (
	"context"
	"learn-gin/internal/model"

	"gorm.io/gorm"
)

type WalletRepo interface {
	Create(ctx context.Context, db *gorm.DB, wallet *model.Wallet) error
}

type walletRepo struct {}

func NewWalletRepo() WalletRepo {
	return &walletRepo{}
}

func (r *walletRepo) Create(ctx context.Context, db *gorm.DB, wallet *model.Wallet) error {
	if err := db.WithContext(ctx).Create(wallet).Error; err != nil {
		if isDuplicateKey(err) {
			return ErrDuplicateKey
		}

		return err
	}
	return nil
}