package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	// Postgres driver
	_ "github.com/lib/pq"
)

type Client struct {
	DB *sqlx.DB
}

type Credentials struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Name     string
	Options  string
}

func New(ctx context.Context, creds Credentials) (*Client, error) {
	c := &Client{}

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		creds.User,
		creds.Password,
		creds.Host,
		creds.Port,
		creds.Name,
		creds.Options,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database, %+v", err)
	}
	c.DB = db

	return c, nil
}

func (c *Client) Close() error {
	if err := c.DB.Close(); err != nil {
		return err
	}

	return nil
}

type Closer interface {
	Close() error
}

func closeAndLogErr(closer Closer) {
	err := closer.Close()
	if err != nil {
		slog.Error("Error closing: %v", err)
	}
}
