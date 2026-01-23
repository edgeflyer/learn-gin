package service

import (
	"context"
	"errors"
	"learn-gin/internal/data"
	"learn-gin/internal/model"
	repo "learn-gin/internal/repository"
	"learn-gin/internal/response"
	"learn-gin/internal/utils"
	"time"
)

type SessionService interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	Refresh(ctx context.Context, refreshTokenPlain string) (string, error)
	Logout(ctx context.Context, refreshTokenPlain string) error
}

type sessionService struct {
	tm data.TransactionManager
	userRepo repo.UserRepo
	rtRepo repo.RefreshTokenRepo
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewSessionService(userRepo repo.UserRepo, rtRepo repo.RefreshTokenRepo, tm data.TransactionManager) SessionService {
	return &sessionService{tm: tm, userRepo: userRepo, rtRepo: rtRepo}
}

func (s *sessionService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	u, err := s.userRepo.GetByUsername(ctx, s.tm.GetDB(), username)
	if err != nil {
		return nil, response.UserInvalid
	}

	if !(utils.CheckPassword(password, u.Password)) {
		return nil, response.UserInvalid
	}

	// jwt
	access, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		return nil,  err
	}
	// refresh token
	plain, err := utils.NewRefreshTokenPlain()
	if err != nil {
		return nil, err
	}

	hash := utils.HashRefreshToken(plain)

	now := time.Now()
	rt := &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: now.Add(14 * 24 * time.Hour),
		CreatedAt: now,
	}

	if err := s.rtRepo.Create(ctx, s.tm.GetDB(), rt); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: access,
		RefreshToken: plain,
	}, nil

}

func (s *sessionService) Refresh(ctx context.Context, refreshTokenPlain string) (string, error) {
	hash := utils.HashRefreshToken(refreshTokenPlain)

	rt, err := s.rtRepo.GetByHash(ctx, s.tm.GetDB(), hash)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return "", response.ErrRefreshTokenInvalid
		}
		return "", err
	}

	if rt.RevokedAt != nil {
		return "", response.ErrRefreshTokenRevoked
	}
	if time.Now().After(rt.ExpiresAt) {
		return "", response.ErrRefreshTokenExpired
	}

	u, err := s.userRepo.GetByID(ctx, s.tm.GetDB(), rt.UserID)
	if err != nil {
		return "", err
	}


	// 生成新的access token
	access, err := utils.GenerateToken(rt.UserID, u.Username)
	if err != nil {
		return "", err
	}

	// 生成新的refresh
	return access, nil
}

func (s *sessionService) Logout(ctx context.Context, refreshTokenPlain string) error {
	hash := utils.HashRefreshToken(refreshTokenPlain)

	rt, err := s.rtRepo.GetByHash(ctx, s.tm.GetDB(), hash)
	if err != nil {
		// 登出做幂等，不存在也成功
		if errors.Is(err, repo.ErrNotFound) {
			return nil
		}
		return err
	}

	now := time.Now()
	return s.rtRepo.RevokeByID(ctx, s.tm.GetDB(), rt.ID, now)
}