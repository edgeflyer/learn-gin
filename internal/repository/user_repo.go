package repo

import (
	"context"
	"database/sql"
	"learn-gin/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, db DBTX, user *model.User) error
	GetByUsername(ctx context.Context, db DBTX, username string) (*model.User, error)
}

type userRepo struct {}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

// 创建用户,使用t.tx
func(r *userRepo) Create(ctx context.Context, db DBTX, user *model.User) error {
	query := `
		insert into users (username, password, created_at, updated_at)
		values (?, ?, ?, ?)
	`
	// 这里可能是db，也可能是tx
	result, err := db.ExecContext(ctx, query, user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *userRepo)GetByUsername(ctx context.Context, db DBTX, username string) (*model.User, error) {
	query := `
		select id, username, password, created_at, updated_at
		from users where username = ?		
	`
	row := db.QueryRowContext(ctx, query, username)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
