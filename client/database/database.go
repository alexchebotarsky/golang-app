package database

import (
	"context"
	"fmt"

	"github.com/goodleby/pure-go-server/config"
	"github.com/jmoiron/sqlx"

	// Postgres driver
	_ "github.com/lib/pq"
)

// Article is a database article.
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	ID          string `json:"id"`
}

// Client is a database client.
type Client struct {
	DB                *sqlx.DB
	getArticlesStmt   *sqlx.Stmt
	getArticleStmt    *sqlx.NamedStmt
	addArticleStmt    *sqlx.NamedStmt
	removeArticleStmt *sqlx.NamedStmt
	updateArticleStmt *sqlx.NamedStmt
	Articles          []Article
}

// New creates a new database client.
func New(ctx context.Context, config *config.Config) (*Client, error) {
	var c Client

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		config.DatabaseOptions,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database, %+v", err)
	}
	c.DB = db

	if err := c.prepareGetArticles(); err != nil {
		return nil, err
	}

	if err := c.prepareGetArticle(); err != nil {
		return nil, err
	}

	if err := c.prepareAddArticle(); err != nil {
		return nil, err
	}

	if err := c.prepareRemoveArticle(); err != nil {
		return nil, err
	}

	if err := c.prepareUpdateArticle(); err != nil {
		return nil, err
	}

	// TODO: remove me
	c.Articles = []Article{
		{Title: "Hello World", Description: "This is a description", Body: "This is the body", ID: "1"},
		{Title: "Hello World 2", Description: "This is a description 2", Body: "This is the body 2", ID: "2"},
	}

	return &c, nil
}

func (c *Client) Close() error {
	if err := c.getArticlesStmt.Close(); err != nil {
		return err
	}

	if err := c.getArticleStmt.Close(); err != nil {
		return err
	}

	if err := c.addArticleStmt.Close(); err != nil {
		return err
	}

	if err := c.removeArticleStmt.Close(); err != nil {
		return err
	}

	if err := c.updateArticleStmt.Close(); err != nil {
		return err
	}

	if err := c.DB.Close(); err != nil {
		return err
	}

	return nil
}
