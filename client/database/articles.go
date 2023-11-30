package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/goodleby/golang-app/article"
	"github.com/goodleby/golang-app/tracing"
	"github.com/jmoiron/sqlx"
)

type ArticleStatements struct {
	SelectAll *sqlx.Stmt
	Select    *sqlx.NamedStmt
	Insert    *sqlx.NamedStmt
	Delete    *sqlx.NamedStmt
	Update    *sqlx.NamedStmt
}

func (s *ArticleStatements) Close() error {
	errStrings := []string{}

	if err := s.Select.Close(); err != nil {
		errStrings = append(errStrings, fmt.Sprintf("error closing select statement: %v", err))
	}

	if len(errStrings) > 0 {
		return errors.New(strings.Join(errStrings, "; "))
	}

	return nil
}

func (c *Client) prepareArticleStatements(ctx context.Context) (*ArticleStatements, error) {
	var statements ArticleStatements
	var err error

	query := `SELECT id, title, description, body FROM articles`
	statements.SelectAll, err = c.DB.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing select all article statement: %v", err)
	}

	query = `SELECT id, title, description, body FROM articles WHERE id = :id`
	statements.Select, err = c.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing select article statement: %v", err)
	}

	query = `INSERT INTO articles (title, description, body)
	        	VALUES (:title, :description, :body)
						RETURNING id, title, description, body`
	statements.Insert, err = c.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing insert article statement: %v", err)
	}

	query = `DELETE FROM articles WHERE id = :id`
	statements.Delete, err = c.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing delete article statement: %v", err)
	}

	query = `UPDATE articles
						SET title = :title, description = :description, body = :body
						WHERE id = :id
						RETURNING id, title, description, body`
	statements.Update, err = c.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing update article statement: %v", err)
	}

	return &statements, nil
}

func (c *Client) SelectAllArticles(ctx context.Context) ([]article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "SelectAllArticles")
	defer span.End()

	var articles []article.Article
	if err := c.ArticleStatements.SelectAll.SelectContext(ctx, &articles); err != nil {
		return nil, fmt.Errorf("error selecting articles: %v", err)
	}

	return articles, nil
}

func (c *Client) SelectArticle(ctx context.Context, id string) (*article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "SelectArticle")
	defer span.End()

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	var article article.Article
	if err := c.ArticleStatements.Select.GetContext(ctx, &article, args); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound{Err: err}
		default:
			return nil, fmt.Errorf("error selecting article with id %q: %v", id, err)
		}
	}

	return &article, nil
}

func (c *Client) InsertArticle(ctx context.Context, payload article.Payload) (*article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "InsertArticle")
	defer span.End()

	args := struct {
		article.Payload
	}{
		Payload: payload,
	}

	var article article.Article
	if err := c.ArticleStatements.Insert.GetContext(ctx, &article, args); err != nil {
		return nil, fmt.Errorf("error inserting an article: %v", err)
	}

	return &article, nil
}

func (c *Client) DeleteArticle(ctx context.Context, id string) error {
	ctx, span := tracing.StartSpan(ctx, "DeleteArticle")
	defer span.End()

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	result, err := c.ArticleStatements.Delete.ExecContext(ctx, args)
	if err != nil {
		return fmt.Errorf("error deleting article: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return ErrNotFound{Err: errors.New("no rows to delete")}
	}

	return nil
}

func (c *Client) UpdateArticle(ctx context.Context, id string, payload article.Payload) (*article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "UpdateArticle")
	defer span.End()

	args := struct {
		article.Payload
		ID string `db:"id"`
	}{
		Payload: payload,
		ID:      id,
	}

	var article article.Article
	if err := c.ArticleStatements.Update.GetContext(ctx, &article, args); err != nil {
		return nil, fmt.Errorf("error updating article: %v", err)
	}

	return &article, nil
}
