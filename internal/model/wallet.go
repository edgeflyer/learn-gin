package model

import "time"

type Wallet struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;uniqueIndex;index"`
	Balance   int64     `json:"balance" gorm:"not null;default:0"` // 推荐 int64 分
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}