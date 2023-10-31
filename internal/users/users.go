package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserNotExists = errors.New("user doesn't exists")
	ErrUserExists    = errors.New("user exists")
)

const (
	// QUERIES
	addQueryString    = `INSERT INTO public.users (id) VALUES ($1)`
	deleteQueryString = `DELETE FROM public.users WHERE id=($1)`

	getQueryString    = `SELECT * FROM public.users WHERE id=($1)`
	getAllQueryString = `SELECT id FROM public.users`
)

type Repository interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type UserRepository struct {
	db Repository
}

func NewUserRepository(repo Repository) *UserRepository {
	return &UserRepository{
		db: repo,
	}
}

func (r *UserRepository) CheckIfUserExists(ctx context.Context, chatID int64) error {
	result := r.db.QueryRow(ctx, getQueryString, chatID)
	err := result.Scan(&chatID)
	switch err {
	case pgx.ErrNoRows:
		return ErrUserNotExists
	case nil:
		return ErrUserExists
	default:
		return err
	}
}

func (r *UserRepository) AddUserToDB(ctx context.Context, chatID int64) error {
	_, err := r.db.Exec(ctx, addQueryString, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUserFromDB(ctx context.Context, chatID int64) error {
	_, err := r.db.Exec(ctx, deleteQueryString, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUsersList(ctx context.Context) ([]int64, error) {
	rows, err := r.db.Query(ctx, getAllQueryString)
	if err != nil {
		return nil, err
	}
	userList := []int64{}
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		userList = append(userList, id)
	}
	return userList, nil
}
