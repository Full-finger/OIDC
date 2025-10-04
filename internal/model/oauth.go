// Package model defines the data structures used in the OIDC application.
package model

import (
	"time"
)

// Client represents an OAuth client.
type Client struct {
	ID                int64     `db:"id"`
	ClientID          string    `db:"client_id"`
	ClientSecretHash  string    `db:"client_secret_hash"`
	Name              string    `db:"name"`
	RedirectURIs      []string  `db:"redirect_uris"`
	Scopes            []string  `db:"scopes"`
	CreatedAt         time.Time `db:"created_at"`
}

// SafeClient returns a safe version of Client without sensitive information.
func (c *Client) SafeClient() SafeClient {
	return SafeClient{
		ID:           c.ID,
		ClientID:     c.ClientID,
		Name:         c.Name,
		RedirectURIs: c.RedirectURIs,
		Scopes:       c.Scopes,
		CreatedAt:    c.CreatedAt,
	}
}

// SafeClient is a safe subset of Client for external exposure.
type SafeClient struct {
	ID           int64     `json:"id"`
	ClientID     string    `json:"client_id"`
	Name         string    `json:"name"`
	RedirectURIs []string  `json:"redirect_uris"`
	Scopes       []string  `json:"scopes"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuthorizationCode represents an OAuth authorization code.
type AuthorizationCode struct {
	Code                string    `db:"code"`
	ClientID            string    `db:"client_id"`
	UserID              int64     `db:"user_id"`
	RedirectURI         string    `db:"redirect_uri"`
	Scopes              []string  `db:"scopes"`
	ExpiresAt           time.Time `db:"expires_at"`
	CodeChallenge       *string   `db:"code_challenge"`
	CodeChallengeMethod *string   `db:"code_challenge_method"`
}

// RefreshToken represents an OAuth refresh token.
type RefreshToken struct {
	ID          int64     `db:"id"`
	TokenHash   string    `db:"token_hash"`
	UserID      int64     `db:"user_id"`
	ClientID    string    `db:"client_id"`
	Scopes      []string  `db:"scopes"`
	ExpiresAt   time.Time `db:"expires_at"`
	RevokedAt   *time.Time `db:"revoked_at"`
}