package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goodleby/golang-app/article"
	"github.com/goodleby/golang-app/tracing"
)

func (c *Client) SelectAllArticles(ctx context.Context) ([]article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "SelectAllArticles")
	defer span.End()

	query := `SELECT id, title, description, body FROM articles`

	var articles []article.Article
	if err := c.DB.SelectContext(ctx, &articles, query); err != nil {
		return nil, fmt.Errorf("error selecting articles: %v", err)
	}

	return articles, nil
}

func (c *Client) SelectArticle(ctx context.Context, id string) (*article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "SelectArticle")
	defer span.End()

	stmt, err := c.DB.PrepareNamedContext(ctx, `SELECT id, title, description, body FROM articles WHERE id = :id`)
	if err != nil {
		return nil, fmt.Errorf("error preparing named statement: %v", err)
	}
	defer closeAndLogErr(stmt)

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	var article article.Article
	if err := stmt.GetContext(ctx, &article, args); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound{Err: err}
		default:
			return nil, fmt.Errorf("error selecting article with id %q: %v", id, err)
		}
	}

	return &article, nil
}

func (c *Client) InsertArticle(ctx context.Context, payload article.Payload) error {
	ctx, span := tracing.StartSpan(ctx, "InsertArticle")
	defer span.End()

	query := `INSERT INTO articles (title, description, body) VALUES (:title, :description, :body)`

	args := struct {
		article.Payload
	}{
		Payload: payload,
	}

	if _, err := c.DB.NamedExecContext(ctx, query, args); err != nil {
		return fmt.Errorf("error inserting an article: %v", err)
	}

	return nil
}

func (c *Client) DeleteArticle(ctx context.Context, id string) error {
	ctx, span := tracing.StartSpan(ctx, "DeleteArticle")
	defer span.End()

	query := `DELETE FROM articles WHERE id = :id`

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	if _, err := c.DB.NamedExecContext(ctx, query, args); err != nil {
		return fmt.Errorf("error deleting article: %v", err)
	}

	return nil
}

func (c *Client) UpdateArticle(ctx context.Context, id string, payload article.Payload) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateArticle")
	defer span.End()

	query := `UPDATE articles SET title = :title, description = :description, body = :body WHERE id = :id`

	args := struct {
		article.Payload
		ID string `db:"id"`
	}{
		Payload: payload,
		ID:      id,
	}

	if _, err := c.DB.NamedExecContext(ctx, query, args); err != nil {
		return fmt.Errorf("error updating article: %v", err)
	}

	return nil
}
