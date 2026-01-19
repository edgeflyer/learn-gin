package service

import (
	"context"
	"database/sql"
	"learn-gin/internal/data"
	"learn-gin/internal/model"
	repo "learn-gin/internal/repository"
	"learn-gin/internal/response"
	"learn-gin/internal/utils"
	"time"
)


type UserService interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (string, error)
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
	status, err := s.userRepo.GetByUsername(ctx, s.tm.GetDB(), username)
	if err != nil {
		return err // 这是一个未知错误，直接抛出
	}
	if status != nil {
		return response.UserExists
	}

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
	return s.tm.ExecTx(ctx, func(ctx context.Context, tx repo.DBTX) error {

		// 将tx显示传给UserRepo
		if err := s.userRepo.Create(ctx, tx, user); err != nil {
			return err
		}

		// 将统一个tx传给WalletRepo
		// 因为用的是同一个tx，所以他俩在同一个数据库事务里
		wallet := &model.Wallet{
			UserID: user.ID,
			Balance: 0,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.walletRepo.Create(ctx, tx, wallet); err != nil {
			return err
		}

		return nil
	})
}

func (s *userService) Login(ctx context.Context,username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, s.tm.GetDB(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", response.UserNotFound
		}
		return "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}