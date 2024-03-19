package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/goodleby/golang-app/article"
	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/tracing"
	"github.com/jmoiron/sqlx"
)

type ArticleStmt struct {
	SelectAll *sqlx.Stmt
	Select    *sqlx.NamedStmt
	Insert    *sqlx.NamedStmt
	Delete    *sqlx.NamedStmt
	Update    *sqlx.NamedStmt
}

func (c *Client) prepareArticleStatements(ctx context.Context) (*ArticleStmt, error) {
	var articleStmt ArticleStmt
	var err error

	articleStmt.SelectAll, err = c.prepareSelectAllArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing select all articles statement: %v", err)
	}

	articleStmt.Select, err = c.prepareSelectArticle(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing select article statement: %v", err)
	}

	articleStmt.Insert, err = c.prepareInsertArticle(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing insert article statement: %v", err)
	}

	articleStmt.Delete, err = c.prepareDeleteArticle(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing delete article statement: %v", err)
	}

	articleStmt.Update, err = c.prepareUpdateArticle(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing update article statement: %v", err)
	}

	return &articleStmt, nil
}

func (articleStmt *ArticleStmt) Close() error {
	errs := []error{}

	err := articleStmt.SelectAll.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing select all articles statement: %v", err))
	}

	err = articleStmt.Select.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing select article statement: %v", err))
	}

	err = articleStmt.Insert.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing insert article statement: %v", err))
	}

	err = articleStmt.Delete.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing delete article statement: %v", err))
	}

	err = articleStmt.Update.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing update article statement: %v", err))
	}

	if len(errs) > 0 {
		return &client.ErrMultiple{Errs: errs}
	}

	return nil
}

func (c *Client) prepareSelectAllArticles(ctx context.Context) (*sqlx.Stmt, error) {
	query := "SELECT id, title, description, body FROM articles"
	return c.DB.PreparexContext(ctx, query)
}

func (c *Client) SelectAllArticles(ctx context.Context) ([]article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "SelectAllArticles")
	defer span.End()

	articles := []article.Article{}
	err := c.ArticleStmt.SelectAll.SelectContext(ctx, &articles)
	if err != nil {
		return nil, fmt.Errorf("error selecting articles: %v", err)
	}

	return articles, nil
}

func (c *Client) prepareSelectArticle(ctx context.Context) (*sqlx.NamedStmt, error) {
	query := "SELECT id, title, description, body FROM articles WHERE id = :id"
	return c.DB.PrepareNamedContext(ctx, query)
}

func (c *Client) SelectArticle(ctx context.Context, id int) (*article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "SelectArticle")
	defer span.End()

	args := struct {
		ID int `db:"id"`
	}{
		ID: id,
	}

	var article article.Article
	err := c.ArticleStmt.Select.GetContext(ctx, &article, args)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, &client.ErrNotFound{Err: fmt.Errorf("article with id %d not found: %v", id, err)}
		default:
			return nil, fmt.Errorf("error selecting article with id %d: %v", id, err)
		}
	}

	return &article, nil
}

func (c *Client) prepareInsertArticle(ctx context.Context) (*sqlx.NamedStmt, error) {
	query := `INSERT INTO articles (title, description, body)
	        	VALUES (:title, :description, :body)
						RETURNING id, title, description, body`
	return c.DB.PrepareNamedContext(ctx, query)
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
	err := c.ArticleStmt.Insert.GetContext(ctx, &article, args)
	if err != nil {
		return nil, fmt.Errorf("error inserting an article: %v", err)
	}

	return &article, nil
}

func (c *Client) prepareDeleteArticle(ctx context.Context) (*sqlx.NamedStmt, error) {
	query := `DELETE FROM articles WHERE id = :id`
	return c.DB.PrepareNamedContext(ctx, query)
}

func (c *Client) DeleteArticle(ctx context.Context, id int) error {
	ctx, span := tracing.StartSpan(ctx, "DeleteArticle")
	defer span.End()

	args := struct {
		ID int `db:"id"`
	}{
		ID: id,
	}

	result, err := c.ArticleStmt.Delete.ExecContext(ctx, args)
	if err != nil {
		return fmt.Errorf("error deleting article with id %d: %v", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return &client.ErrNotFound{Err: errors.New("no rows to delete")}
	}

	return nil
}

func (c *Client) prepareUpdateArticle(ctx context.Context) (*sqlx.NamedStmt, error) {
	query := `UPDATE articles
						SET title = :title, description = :description, body = :body
						WHERE id = :id
						RETURNING id, title, description, body`
	return c.DB.PrepareNamedContext(ctx, query)
}

func (c *Client) UpdateArticle(ctx context.Context, id int, payload article.Payload) (*article.Article, error) {
	ctx, span := tracing.StartSpan(ctx, "UpdateArticle")
	defer span.End()

	args := struct {
		article.Payload
		ID int `db:"id"`
	}{
		Payload: payload,
		ID:      id,
	}

	var article article.Article
	err := c.ArticleStmt.Update.GetContext(ctx, &article, args)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, &client.ErrNotFound{Err: fmt.Errorf("article with id %d not found: %v", id, err)}
		default:
			return nil, fmt.Errorf("error updating article with id %d: %v", id, err)
		}
	}

	return &article, nil
}
