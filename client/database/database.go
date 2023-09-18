package database

import (
	"context"
	"fmt"

	"github.com/goodleby/golang-server/env"
	"github.com/jmoiron/sqlx"

	// Postgres driver
	_ "github.com/lib/pq"
)

type Client struct {
	DB *sqlx.DB
}

func New(ctx context.Context, env *env.Config) (*Client, error) {
	var c Client

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseHost,
		env.DatabasePort,
		env.DatabaseName,
		env.DatabaseOptions,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database, %+v", err)
	}
	c.DB = db

	return &c, nil
}

func (c *Client) Close() error {
	if err := c.DB.Close(); err != nil {
		return err
	}

	return nil
}
