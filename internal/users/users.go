package users

import (
	"context"
	"errors"
)

var (
	ErrUserNotExists = errors.New("user doesn't exists")
	ErrUserExists    = errors.New("user exists")
)

type Repository interface {
	CheckIfUserExists(ctx context.Context, chatID int64) error
	AddUserToDB(ctx context.Context, chatID int64) error
	DeleteUserFromDB(ctx context.Context, chatID int64) error
	GetUsersList(ctx context.Context) ([]int64, error)
}

type UserRepository struct {
	db Repository
}

func (r *UserRepository) CheckIfUserExists(ctx context.Context, chatID int64) error {
	return r.db.CheckIfUserExists(ctx, chatID)
}

func (r *UserRepository) AddUserToDB(ctx context.Context, chatID int64) error {
	return r.db.AddUserToDB(ctx, chatID)
}

func (r *UserRepository) DeleteUserFromDB(ctx context.Context, chatID int64) error {
	return r.db.DeleteUserFromDB(ctx, chatID)
}

func (r *UserRepository) GetUsersList(ctx context.Context) ([]int64, error) {
	return r.db.GetUsersList(ctx)
}
