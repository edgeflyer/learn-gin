package model

import "time"

type User struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"size:64;not null;uniqueIndex"`
	Password string `json:"-" gorm:"column:password_hash;size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}