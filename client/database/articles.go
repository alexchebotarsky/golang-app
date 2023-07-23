package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/goodleby/golang-server/tracing"
)

func (c *Client) prepareGetArticles() error {
	stmt, err := c.DB.Preparex(`SELECT id, title, description, body FROM articles`)

	if err != nil {
		return fmt.Errorf("error preparing get articles statement: %v", err)
	}

	c.getArticlesStmt = stmt

	return nil
}

// FetchAllArticles fetches all articles.
func (c *Client) FetchAllArticles(ctx context.Context) ([]Article, error) {
	ctx, span := tracing.Span(ctx, "FetchAllArticles")
	defer span.End()

	rows, err := c.getArticlesStmt.QueryContext(ctx)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, &ErrNotFound{Err: err}
		default:
			return nil, fmt.Errorf("error querying articles: %v", err)
		}
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing database rows: %v", err)
		}
	}()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Description,
			&article.Body,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (c *Client) prepareGetArticle() error {
	stmt, err := c.DB.PrepareNamed(`SELECT id, title, description, body FROM articles WHERE id = :id`)

	if err != nil {
		return fmt.Errorf("error preparing get articles statement: %v", err)
	}

	c.getArticleStmt = stmt

	return nil
}

// FetchArticle fetches an article by id.
func (c *Client) FetchArticle(ctx context.Context, id string) (*Article, error) {
	ctx, span := tracing.Span(ctx, "FetchArticle")
	defer span.End()

	var article Article

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	err := c.getArticleStmt.QueryRowContext(ctx, args).Scan(&article.ID, &article.Title, &article.Description, &article.Body)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, &ErrNotFound{Err: err}
		default:
			return nil, fmt.Errorf("error querying article with id %q: %v", id, err)
		}
	}

	return &article, nil
}

func (c *Client) prepareAddArticle() error {
	stmt, err := c.DB.PrepareNamed(`INSERT INTO articles (id, title, description, body) VALUES (:id, :title, :description, :body)`)

	if err != nil {
		return fmt.Errorf("error preparing add article statement: %v", err)
	}

	c.addArticleStmt = stmt

	return nil
}

// AddArticle adds an article.
func (c *Client) AddArticle(ctx context.Context, article Article) error {
	ctx, span := tracing.Span(ctx, "AddArticle")
	defer span.End()

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

	_, err := c.addArticleStmt.ExecContext(ctx, args)
	if err != nil {
		return fmt.Errorf("error adding an article: %v", err)
	}

	return nil
}

func (c *Client) prepareRemoveArticle() error {
	stmt, err := c.DB.PrepareNamed(`DELETE FROM articles WHERE id = :id`)
	if err != nil {
		return fmt.Errorf("error preparing remove article statement: %v", err)
	}

	c.removeArticleStmt = stmt

	return nil
}

// RemoveArticle removes an article.
func (c *Client) RemoveArticle(ctx context.Context, id string) error {
	ctx, span := tracing.Span(ctx, "RemoveArticle")
	defer span.End()

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	_, err := c.removeArticleStmt.ExecContext(ctx, args)
	if err != nil {
		return fmt.Errorf("error removing article: %v", err)
	}

	return nil
}

func (c *Client) prepareUpdateArticle() error {
	stmt, err := c.DB.PrepareNamed(`UPDATE articles SET id = :new_id, title = :new_title, description = :new_description, body = :new_body WHERE id = :id`)
	if err != nil {
		return fmt.Errorf("error preparing update article statement: %v", err)
	}

	c.updateArticleStmt = stmt

	return nil
}

// UpdateArticle updates an article.
func (c *Client) UpdateArticle(ctx context.Context, id string, article Article) error {
	ctx, span := tracing.Span(ctx, "UpdateArticle")
	defer span.End()

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

	_, err := c.updateArticleStmt.ExecContext(ctx, args)
	if err != nil {
		return fmt.Errorf("error updating article: %v", err)
	}

	return nil
}
