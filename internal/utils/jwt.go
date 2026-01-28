package utils

import (
	"errors"
	"learn-gin/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// GenerateToken 生成 Token
func GenerateToken(userID int64, username string) (string, error) {
    now := time.Now()
    expireTime := now.Add(time.Duration(config.Conf.JWT.Expire) * time.Minute)

    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expireTime),
            IssuedAt:  jwt.NewNumericDate(now),
            Issuer:    "phase4", // 签发人
        },
    }

    tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // 
    return tokenClaims.SignedString([]byte(config.Conf.JWT.Secret))
}

// ParseToken 解析 Token
func ParseToken(token string) (*Claims, error) {
    tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {

        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(config.Conf.JWT.Secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}