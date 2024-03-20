package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/goodleby/golang-app/client"
	"github.com/jmoiron/sqlx"

	// Postgres driver
	_ "github.com/lib/pq"
)

type Client struct {
	DB          *sqlx.DB
	ArticleStmt *ArticleStmt
}

func New(ctx context.Context, creds Credentials) (*Client, error) {
	var c Client
	var err error

	c.DB, err = sqlx.ConnectContext(ctx, "postgres", creds.ToConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	c.ArticleStmt, err = c.prepareArticleStatements(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing article statements: %v", err)
	}

	return &c, nil
}

type Credentials struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Name     string
	Options  string
}

func (c *Credentials) ToConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		c.User,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Name,
		c.Options,
	)
}

func (c *Client) Close() error {
	errs := []error{}

	err := c.ArticleStmt.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing article statements: %v", err))
	}

	err = c.DB.Close()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return &client.ErrMultiple{Errs: errs}
	}

	return nil
}
