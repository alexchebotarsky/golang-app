package database

import (
	"context"
	"fmt"

	"github.com/goodleby/golang-server/config"
	"github.com/jmoiron/sqlx"

	// Postgres driver
	_ "github.com/lib/pq"
)

type Client struct {
	DB *sqlx.DB
}

func New(ctx context.Context, config *config.Config) (*Client, error) {
	var c Client

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		config.DatabaseOptions,
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
