package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/goodleby/golang-server/tracing"
)

func (c *Client) SelectAllArticles(ctx context.Context) ([]Article, error) {
	ctx, span := tracing.Span(ctx, "SelectAllArticles")
	defer span.End()

	query := `SELECT id, title, description, body FROM articles`

	var articles []Article
	if err := c.DB.SelectContext(ctx, &articles, query); err != nil {
		return nil, fmt.Errorf("error selecting articles: %v", err)
	}

	return articles, nil
}

func (c *Client) SelectArticle(ctx context.Context, id string) (*Article, error) {
	ctx, span := tracing.Span(ctx, "SelectArticle")
	defer span.End()

	stmt, err := c.DB.PrepareNamedContext(ctx, `SELECT id, title, description, body FROM articles WHERE id = :id`)
	if err != nil {
		return nil, fmt.Errorf("error preparing named statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing prepared named statement: %v", err)
		}
	}()

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	var article Article
	if err := stmt.GetContext(ctx, &article, args); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, &ErrNotFound{Err: err}
		default:
			return nil, fmt.Errorf("error getting article with id %q: %v", id, err)
		}
	}

	return &article, nil
}

func (c *Client) AddArticle(ctx context.Context, article Article) error {
	ctx, span := tracing.Span(ctx, "AddArticle")
	defer span.End()

	query := `INSERT INTO articles (id, title, description, body) VALUES (:id, :title, :description, :body)`

	args := struct {
		ID          string `db:"id"`
		Title       string `db:"title"`
		Description string `db:"description"`
		Body        string `db:"body"`
	}{
		ID:          article.ID,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
	}

	if _, err := c.DB.NamedExecContext(ctx, query, args); err != nil {
		return fmt.Errorf("error adding an article: %v", err)
	}

	return nil
}

func (c *Client) RemoveArticle(ctx context.Context, id string) error {
	ctx, span := tracing.Span(ctx, "RemoveArticle")
	defer span.End()

	query := `DELETE FROM articles WHERE id = :id`

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	_, err := c.DB.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error removing article: %v", err)
	}

	return nil
}

func (c *Client) UpdateArticle(ctx context.Context, id string, article Article) error {
	ctx, span := tracing.Span(ctx, "UpdateArticle")
	defer span.End()

	query := `UPDATE articles SET id = :new_id, title = :new_title, description = :new_description, body = :new_body WHERE id = :id`

	args := struct {
		ID             string `db:"id"`
		NewID          string `db:"new_id"`
		NewTitle       string `db:"new_title"`
		NewDescription string `db:"new_description"`
		NewBody        string `db:"new_body"`
	}{
		ID:             id,
		NewID:          article.ID,
		NewTitle:       article.Title,
		NewDescription: article.Description,
		NewBody:        article.Body,
	}

	if _, err := c.DB.NamedExecContext(ctx, query, args); err != nil {
		return fmt.Errorf("error updating article: %v", err)
	}

	return nil
}
