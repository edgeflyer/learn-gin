package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// DefaultCost 目前是10，是一个在性能和安全性之间很好的平衡
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}