package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Client struct {
	conn *pgx.Conn
}

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func New(ctx context.Context, conf Config) (*Client, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY);`)

	if err != nil {
		return nil, fmt.Errorf("can't apply schema: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

func (c *Client) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.conn.QueryRow(ctx, sql, args...)
}

func (c *Client) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.conn.Query(ctx, sql, args...)
}

func (c *Client) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.conn.Exec(ctx, sql, args...)
}
