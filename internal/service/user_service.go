package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"learn-gin/internal/data"
	"learn-gin/internal/logger"
	"learn-gin/internal/model"
	"learn-gin/internal/redis"
	repo "learn-gin/internal/repository"
	"learn-gin/internal/response"
	"learn-gin/internal/utils"
	"time"

	"gorm.io/gorm"
)


type UserService interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (string, error)
	GetUserProfile(ctx context.Context, userID int64) (*model.User, error)
}

type userService struct {
	tm data.TransactionManager // 事务管理器
	userRepo repo.UserRepo
	walletRepo repo.WalletRepo
}

func NewUserService(tm data.TransactionManager, uRepo repo.UserRepo, wRepo repo.WalletRepo) UserService {
	return &userService{
		tm: tm,
		userRepo: uRepo,
		walletRepo: wRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, password string) error {

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &model.User{
		Username: username,
		Password: hashedPwd,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 开启事物
	return s.tm.ExecTx(ctx, func(ctx context.Context, tx *gorm.DB) error {

		// 将tx显示传给UserRepo
		if err := s.userRepo.Create(ctx, tx, user); err != nil {
			if errors.Is(err, repo.ErrDuplicateKey) {
				return response.UserExists
			}
			return err
		}

		// 将同一个tx传给WalletRepo
		// 因为用的是同一个tx，所以他俩在同一个数据库事务里
		wallet := &model.Wallet{
			UserID: user.ID,
			Balance: 0,
			CreatedAt: now,
			UpdatedAt: now,
		}
		// 按照道理来说，应该不会重复吧？
		if err := s.walletRepo.Create(ctx, tx, wallet); err != nil {
			return err
		}

		return nil
	})
}

func (s *userService) Login(ctx context.Context,username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, s.tm.GetDB(), username)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return "", response.UserNotFound
		}
		return "", err
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", response.UserInvalid
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetUserProfile(ctx context.Context, userID int64) (*model.User, error) {

	// 定义缓存
	cacheKey := fmt.Sprintf("cache:user:profile:%d", userID)

	// 尝试从redis中读取数据
	val, err := redis.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存中有数据
		var user model.User
		// 反序列化json到struct
		if jsonErr := json.Unmarshal([]byte(val), &user); jsonErr == nil {
			logger.Log.Info("[cache hit]命中缓存")
			return &user, nil
		}
	}
	// 没有命中，查询数据库
	user, err := s.userRepo.GetByID(ctx, s.tm.GetDB(), userID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, response.UserNotFound
		}
		return nil, err
	}

	data, _ := json.Marshal(user)
	redis.RDB.Set(ctx, cacheKey, data, 30*time.Minute)
	return user, nil
}