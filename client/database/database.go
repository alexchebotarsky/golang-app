package database

import "fmt"

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Id          string `json:"id"`
}

type Client struct {
	Articles []Article
}

func New() (*Client, error) {
	var c Client

	c.Articles = []Article{
		{Title: "Hello World", Description: "This is a description", Body: "This is the body", Id: "1"},
		{Title: "Hello World 2", Description: "This is a description 2", Body: "This is the body 2", Id: "2"},
	}

	return &c, nil
}

func (c *Client) FetchArticles() ([]Article, error) {
	return c.Articles, nil
}

func (c *Client) FetchArticle(id string) (Article, error) {
	for _, article := range c.Articles {
		if article.Id == id {
			return article, nil
		}
	}

	return Article{}, fmt.Errorf("Article with id '%s' not found", id)
}

func (c *Client) CreateArticle(article Article) (Article, error) {
	c.Articles = append(c.Articles, article)

	return article, nil
}

func (c *Client) UpdateArticle(id string, article Article) (Article, error) {
	for i, a := range c.Articles {
		if a.Id == id {
			c.Articles[i] = article
			return article, nil
		}
	}

	return Article{}, fmt.Errorf("Article with id '%s' not found", id)
}

func (c *Client) DeleteArticle(id string) error {
	for i, article := range c.Articles {
		if article.Id == id {
			c.Articles = append(c.Articles[:i], c.Articles[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Article with id '%s' not found", id)
}
