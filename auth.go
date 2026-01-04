package gopholler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// getValidToken returns a valid OAuth2 token, refreshing if necessary.
func (c *Client) getValidToken(ctx context.Context) (*OAuth2Token, error) {
	c.tokenMu.RLock()
	if c.token != nil && time.Now().Before(c.tokenExpiry) {
		token := c.token
		c.tokenMu.RUnlock()
		return token, nil
	}
	c.tokenMu.RUnlock()

	return c.refreshToken(ctx)
}

// refreshToken obtains a new OAuth2 token.
func (c *Client) refreshToken(ctx context.Context) (*OAuth2Token, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// Double-check after acquiring write lock
	if c.token != nil && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	token, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	c.token = token

	// Parse expiry time (comes as a string like "3599")
	expiresIn, err := strconv.Atoi(token.ExpiresIn)
	if err != nil {
		expiresIn = 3600 // Default to 1 hour
	}
	// Subtract 60 seconds to refresh before actual expiry
	c.tokenExpiry = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

	return c.token, nil
}

// authenticate performs the OAuth2 client credentials flow.
func (c *Client) authenticate(ctx context.Context) (*OAuth2Token, error) {
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", strings.Join(c.scopes, " "))

	reqURL := c.oauthURL + "/oauth/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Code:       "AUTH_FAILED",
			Issue:      string(body),
			StatusCode: resp.StatusCode,
			Err:        ErrAuthentication,
		}
	}

	var token OAuth2Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	return &token, nil
}

// InvalidateToken clears the cached token, forcing a refresh on the next request.
func (c *Client) InvalidateToken() {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token = nil
	c.tokenExpiry = time.Time{}
}
