package repo

import (
	"context"
	"errors"
	"learn-gin/internal/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, db *gorm.DB, user *model.User) error
	GetByUsername(ctx context.Context, db *gorm.DB, username string) (*model.User, error)
	GetByID(ctx context.Context, db *gorm.DB, id int64) (*model.User, error)
}

type userRepo struct {}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

// 创建用户,使用t.tx
func(r *userRepo) Create(ctx context.Context, db *gorm.DB, user *model.User) error {
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateKey(err) {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}

func (r *userRepo)GetByUsername(ctx context.Context, db *gorm.DB, username string) (*model.User, error) {
	var user model.User

	if err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil

}

func (r *userRepo) GetByID(ctx context.Context, db *gorm.DB, id int64) (*model.User, error) {
	var user model.User

	if err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}