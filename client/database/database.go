package database

import "fmt"

// Article is a database article.
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	ID          string `json:"id"`
}

// Client is a database client.
type Client struct {
	Articles []Article
}

// New creates a new database client.
func New() (*Client, error) {
	var c Client

	c.Articles = []Article{
		{Title: "Hello World", Description: "This is a description", Body: "This is the body", ID: "1"},
		{Title: "Hello World 2", Description: "This is a description 2", Body: "This is the body 2", ID: "2"},
	}

	return &c, nil
}

// FetchAllArticles fetches all articles.
func (c *Client) FetchAllArticles() ([]Article, error) {
	return c.Articles, nil
}

// FetchArticle fetches an article by id.
func (c *Client) FetchArticle(id string) (*Article, error) {
	for _, article := range c.Articles {
		if article.ID == id {
			return &article, nil
		}
	}

	return nil, fmt.Errorf("Article with id '%s' not found", id)
}

// CreateArticle creates an article.
func (c *Client) CreateArticle(article Article) error {
	c.Articles = append(c.Articles, article)

	return nil
}

// UpdateArticle updates an article.
func (c *Client) UpdateArticle(id string, article Article) error {
	for i, a := range c.Articles {
		if a.ID == id {
			c.Articles[i] = article
			return nil
		}
	}

	return fmt.Errorf("Article with id '%s' not found", id)
}

// RemoveArticle removes an article.
func (c *Client) RemoveArticle(id string) error {
	for i, article := range c.Articles {
		if article.ID == id {
			c.Articles = append(c.Articles[:i], c.Articles[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Article with id '%s' not found", id)
}
