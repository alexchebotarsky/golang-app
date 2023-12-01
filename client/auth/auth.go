package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/goodleby/golang-app/client"
	"github.com/goodleby/golang-app/tracing"
)

type Client struct {
	roles         []Role
	authSecret    []byte
	SigningMethod jwt.SigningMethod
}

type Keys struct {
	Admin  string
	Editor string
	Viewer string
}

func New(ctx context.Context, secret string, keys Keys) (*Client, error) {
	c := &Client{}

	c.authSecret = []byte(secret)
	c.roles = []Role{
		{
			Name:        AdminRole,
			AccessLevel: AdminAccess,
			Key:         keys.Admin,
		},
		{
			Name:        EditorRole,
			AccessLevel: EditorAccess,
			Key:         keys.Editor,
		},
		{
			Name:        ViewerRole,
			AccessLevel: ViewerAccess,
			Key:         keys.Viewer,
		},
	}
	c.SigningMethod = jwt.SigningMethodHS256

	return c, nil
}

type Claims struct {
	RoleName    string `json:"roleName"`
	AccessLevel int    `json:"accessLevel"`
	jwt.RegisteredClaims
}

func (c *Client) NewToken(ctx context.Context, roleName, roleKey string) (string, time.Time, error) {
	_, span := tracing.StartSpan(ctx, "NewToken")
	defer span.End()

	role, err := c.validateRole(roleName, roleKey)
	if err != nil {
		return "", time.Time{}, client.ErrUnauthorized{Err: err}
	}

	expires := time.Now().Add(TokenTTL)

	claims := &Claims{
		RoleName:    role.Name,
		AccessLevel: role.AccessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token := jwt.NewWithClaims(c.SigningMethod, claims)

	signedToken, err := token.SignedString(c.authSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating jwt token: %v", err)
	}

	return signedToken, expires, nil
}

func (c *Client) ParseToken(ctx context.Context, tokenString string) (*Claims, error) {
	_, span := tracing.StartSpan(ctx, "ParseToken")
	defer span.End()

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return c.authSecret, nil
		},
		jwt.WithValidMethods([]string{c.SigningMethod.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing auth token: %v", err)
	}

	if !token.Valid {
		return nil, client.ErrUnauthorized{Err: errors.New("invalid auth token")}
	}

	return claims, nil
}

func (c *Client) RefreshToken(ctx context.Context, tokenString string) (string, time.Time, error) {
	_, span := tracing.StartSpan(ctx, "RefreshToken")
	defer span.End()

	claims := &Claims{}

	oldToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return c.authSecret, nil
		},
		jwt.WithValidMethods([]string{c.SigningMethod.Alg()}),
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error parsing auth token: %v", err)
	}

	if !oldToken.Valid {
		return "", time.Time{}, client.ErrUnauthorized{Err: errors.New("invalid auth token")}
	}

	expires := time.Now().Add(TokenTTL)

	claims.ExpiresAt = jwt.NewNumericDate(expires)

	token := jwt.NewWithClaims(c.SigningMethod, claims)

	signedToken, err := token.SignedString(c.authSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating jwt token: %v", err)
	}

	return signedToken, expires, nil
}

type Role struct {
	Name        string
	AccessLevel int
	Key         string
}

func (c *Client) validateRole(roleName, roleKey string) (*Role, error) {
	for _, role := range c.roles {
		if role.Name == roleName && role.Key == roleKey {
			return &role, nil
		}
	}

	return nil, errors.New("invalid role name or role key")
}

const AdminRole = "admin"
const AdminAccess = 30

const EditorRole = "editor"
const EditorAccess = 20

const ViewerRole = "viewer"
const ViewerAccess = 10

const TokenTTL = 1 * time.Hour
