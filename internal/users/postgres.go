package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	// QUERIES
	addQueryString    = `INSERT INTO public.users (id) VALUES ($1)`
	deleteQueryString = `DELETE FROM public.users WHERE id=($1)`

	getQueryString    = `SELECT * FROM public.users WHERE id=($1)`
	getAllQueryString = `SELECT id FROM public.users`
)

type PgxConn interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type PostgresRepository struct {
	db PgxConn
}

func NewPostgresRepository(db PgxConn) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CheckIfUserExists(ctx context.Context, chatID int64) error {
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

func (r *PostgresRepository) AddUserToDB(ctx context.Context, chatID int64) error {
	_, err := r.db.Exec(ctx, addQueryString, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeleteUserFromDB(ctx context.Context, chatID int64) error {
	_, err := r.db.Exec(ctx, deleteQueryString, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetUsersList(ctx context.Context) ([]int64, error) {
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
