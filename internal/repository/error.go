package repo

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrDuplicateKey = errors.New("repo: duplicate key")
	ErrNotFound = errors.New("repo: not found")
)