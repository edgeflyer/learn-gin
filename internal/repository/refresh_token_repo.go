package repo

import (
	"context"
	"errors"
	"learn-gin/internal/model"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepo interface{
	Create(ctx context.Context, db *gorm.DB, rt *model.RefreshToken) error
	GetByHash(ctx context.Context, db *gorm.DB, tokenHash string) (*model.RefreshToken, error)
	RevokeByID(ctx context.Context, db *gorm.DB, id int64, now time.Time) error
}

type refreshTokenRepo struct{}

func NewRefreshTokenRepo() RefreshTokenRepo {
	return &refreshTokenRepo{}
}

func (r *refreshTokenRepo) Create(ctx context.Context, db *gorm.DB, rt *model.RefreshToken) error {
	return db.WithContext(ctx).Create(rt).Error
}

func (r *refreshTokenRepo) GetByHash(ctx context.Context, db *gorm.DB, tokenHash string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &rt, nil
}

func (r *refreshTokenRepo) RevokeByID(ctx context.Context, db *gorm.DB, id int64, now time.Time) error {
	return db.WithContext(ctx).Model(&model.RefreshToken{}).Where("id = ? and revoked_at is null", id).Update("revoked_at", now).Error
}