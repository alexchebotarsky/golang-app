package database

import (
	"context"
	"fmt"

	"github.com/goodleby/golang-app/client"
	"github.com/jmoiron/sqlx"

	// Postgres driver
	_ "github.com/lib/pq"
)

type Client struct {
	DB                *sqlx.DB
	ArticleStatements *ArticleStatements
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
	var err error

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		creds.User,
		creds.Password,
		creds.Host,
		creds.Port,
		creds.Name,
		creds.Options,
	)

	c.DB, err = sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database, %+v", err)
	}

	c.ArticleStatements, err = c.prepareArticleStatements(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing article statements: %v", err)
	}

	return c, nil
}

func (c *Client) Close() error {
	errs := []error{}

	if err := c.ArticleStatements.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing article statements: %v", err))
	}

	if err := c.DB.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return client.ErrMultiple{Errs: errs}
	}

	return nil
}
