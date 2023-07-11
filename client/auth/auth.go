package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/goodleby/golang-server/config"
)

// Client is an auth client.
type Client struct {
	roles         []Role
	authSecret    []byte
	SigningMethod jwt.SigningMethod
}

// New creates a new auth client.
func New(ctx context.Context, config *config.Config) (*Client, error) {
	var c Client

	c.authSecret = []byte(config.AuthSecret)
	c.roles = []Role{
		{
			Name:        AdminRole,
			AccessLevel: AdminAccess,
			Key:         config.AuthAdminKey,
		},
		{
			Name:        EditorRole,
			AccessLevel: EditorAccess,
			Key:         config.AuthEditorKey,
		},
		{
			Name:        ViewerRole,
			AccessLevel: ViewerAccess,
			Key:         config.AuthViewerKey,
		},
	}
	c.SigningMethod = jwt.SigningMethodHS256

	return &c, nil
}

// Claims contains all standard and custom fields of auth JWT.
type Claims struct {
	RoleName    string `json:"roleName"`
	AccessLevel int    `json:"accessLevel"`
	jwt.RegisteredClaims
}

// NewToken creates a new signed JWT provided role credentials.
func (c *Client) NewToken(roleName, roleKey string) (string, time.Time, error) {
	role, err := c.ValidateRole(roleName, roleKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error validating role ket: %v", err)
	}

	expires := time.Now().Add(5 * time.Minute)

	claims := Claims{
		RoleName:    role.Name,
		AccessLevel: role.AccessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token := jwt.NewWithClaims(c.SigningMethod, &claims)

	signedToken, err := token.SignedString(c.authSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating jwt token: %v", err)
	}

	return signedToken, expires, nil
}

// ParseToken parses and validates provided JWT string and checks its access level.
func (c *Client) ParseToken(tokenString string) (*Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return c.authSecret, nil
		},
		jwt.WithValidMethods([]string{c.SigningMethod.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing auth token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid auth token")
	}

	return &claims, nil
}

// RefreshToken creates a new signed JWT provided old JWT.
func (c *Client) RefreshToken(tokenString string) (string, time.Time, error) {
	claims := Claims{}

	oldToken, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return c.authSecret, nil
		},
		jwt.WithValidMethods([]string{c.SigningMethod.Alg()}),
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error parsing auth token: %v", err)
	}

	if !oldToken.Valid {
		return "", time.Time{}, errors.New("invalid auth token")
	}

	expires := time.Now().Add(5 * time.Minute)

	claims.ExpiresAt = jwt.NewNumericDate(expires)

	token := jwt.NewWithClaims(c.SigningMethod, &claims)

	signedToken, err := token.SignedString(c.authSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error creating jwt token: %v", err)
	}

	return signedToken, expires, nil
}

// Role is an authorized role.
type Role struct {
	Name        string
	AccessLevel int
	Key         string
}

// ValidateRole checks role credentials and if valid returns matched role.
func (c *Client) ValidateRole(roleName, roleKey string) (*Role, error) {
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

const AnyAccess = 0
