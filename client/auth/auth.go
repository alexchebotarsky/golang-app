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
	authSecret    []byte
	tokenTTL      time.Duration
	roles         []Role
	SigningMethod jwt.SigningMethod
}

type Keys struct {
	Admin  string
	Editor string
	Viewer string
}

func New(ctx context.Context, secret string, tokenTTL time.Duration, keys Keys) *Client {
	var c Client

	c.authSecret = []byte(secret)
	c.tokenTTL = tokenTTL
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

	return &c
}

type Claims struct {
	RoleName    string `json:"roleName"`
	AccessLevel int    `json:"accessLevel"`
	jwt.RegisteredClaims
}

func (c *Client) newToken(ctx context.Context, claims jwt.Claims) (string, error) {
	_, span := tracing.StartSpan(ctx, "newToken")
	defer span.End()

	token := jwt.NewWithClaims(c.SigningMethod, claims)

	signedToken, err := token.SignedString(c.authSecret)
	if err != nil {
		return "", fmt.Errorf("error signing jwt token: %v", err)
	}

	return signedToken, nil
}

func (c *Client) parseTokenClaims(ctx context.Context, tokenString string) (Claims, error) {
	_, span := tracing.StartSpan(ctx, "parseTokenClaims")
	defer span.End()

	var claims Claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return c.authSecret, nil
		},
		jwt.WithValidMethods([]string{c.SigningMethod.Alg()}),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing jwt token with claims: %v", err)
	}

	if !token.Valid {
		return Claims{}, &client.ErrUnauthorized{Err: errors.New("invalid auth token")}
	}

	return claims, nil
}

func (c *Client) CreateRoleToken(ctx context.Context, roleName, roleKey string) (string, time.Time, error) {
	_, span := tracing.StartSpan(ctx, "CreateRoleToken")
	defer span.End()

	role, err := c.findRole(roleName, roleKey)
	if err != nil {
		return "", time.Time{}, &client.ErrUnauthorized{Err: err}
	}

	expires := time.Now().Add(c.tokenTTL)

	claims := Claims{
		RoleName:    role.Name,
		AccessLevel: role.AccessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token, err := c.newToken(ctx, claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating new token: %v", err)
	}

	return token, expires, nil
}

func (c *Client) RefreshToken(ctx context.Context, tokenString string) (string, time.Time, error) {
	ctx, span := tracing.StartSpan(ctx, "RefreshToken")
	defer span.End()

	claims, err := c.parseTokenClaims(ctx, tokenString)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error parsing token string: %v", err)
	}

	expires := time.Now().Add(c.tokenTTL)

	claims.ExpiresAt = jwt.NewNumericDate(expires)

	token, err := c.newToken(ctx, claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating new token: %v", err)
	}

	return token, expires, nil
}

func (c *Client) CheckTokenAccess(ctx context.Context, tokenString string, expectedAccessLevel int) error {
	ctx, span := tracing.StartSpan(ctx, "CheckTokenAccess")
	defer span.End()

	claims, err := c.parseTokenClaims(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("error parsing token string: %v", err)
	}

	if expectedAccessLevel > claims.AccessLevel {
		return &client.ErrForbidden{Err: errors.New("insufficient access level")}
	}

	return nil
}

type Role struct {
	Name        string
	AccessLevel int
	Key         string
}

func (c *Client) findRole(roleName, roleKey string) (Role, error) {
	for _, role := range c.roles {
		if role.Name == roleName && role.Key == roleKey {
			return role, nil
		}
	}

	return Role{}, errors.New("invalid role name or role key")
}

const AdminRole = "admin"
const AdminAccess = 30

const EditorRole = "editor"
const EditorAccess = 20

const ViewerRole = "viewer"
const ViewerAccess = 10
