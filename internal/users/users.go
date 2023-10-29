package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	UserNotExistsErr = errors.New("user doesn't exists")
	UserExistsErr    = errors.New("user exists")
)

const (
	//QUERIES
	ADD_QUERY_STRING    = `INSERT INTO public.users (id) VALUES ($1)`
	DELETE_QUERY_STRING = `DELETE FROM public.users WHERE id=($1)`

	GET_QUERY_STRING     = `SELECT * FROM public.users WHERE id=($1)`
	GET_ALL_QUERY_STRING = `SELECT id FROM public.users`
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
	result := r.db.QueryRow(ctx, GET_QUERY_STRING, chatID)
	err := result.Scan(&chatID)
	switch err {
	case pgx.ErrNoRows:
		return UserNotExistsErr
	case nil:
		return UserExistsErr
	default:
		return err
	}
}

func (r *UserRepository) AddUserToDB(ctx context.Context, chatID int64) error {
	_, err := r.db.Exec(ctx, ADD_QUERY_STRING, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUserFromDB(ctx context.Context, chatID int64) error {
	_, err := r.db.Exec(ctx, DELETE_QUERY_STRING, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUsersList(ctx context.Context) ([]int64, error) {
	rows, err := r.db.Query(ctx, GET_ALL_QUERY_STRING)
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
