package gopholler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	DefaultBaseURL  = "https://products.api.telstra.com/messaging/v3"
	DefaultOAuthURL = "https://products.api.telstra.com/v2"
	DefaultTimeout  = 30 * time.Second
)

// Client is the Telstra Messaging API client.
type Client struct {
	baseURL    string
	oauthURL   string
	httpClient *http.Client

	clientID     string
	clientSecret string
	scopes       []string

	token       *OAuth2Token
	tokenExpiry time.Time
	tokenMu     sync.RWMutex
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithOAuthURL sets a custom OAuth URL.
func WithOAuthURL(url string) ClientOption {
	return func(c *Client) {
		c.oauthURL = strings.TrimSuffix(url, "/")
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithScopes sets custom OAuth scopes.
func WithScopes(scopes []string) ClientOption {
	return func(c *Client) {
		c.scopes = scopes
	}
}

// NewClient creates a new Telstra Messaging API client.
func NewClient(clientID, clientSecret string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:      DefaultBaseURL,
		oauthURL:     DefaultOAuthURL,
		httpClient:   &http.Client{Timeout: DefaultTimeout},
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes: []string{
			"free-trial-numbers:read",
			"free-trial-numbers:write",
			"messages:read",
			"messages:write",
			"reports:read",
			"reports:write",
			"virtual-numbers:read",
			"virtual-numbers:write",
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// request performs an HTTP request with automatic token refresh.
func (c *Client) request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	token, err := c.getValidToken(ctx)
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Language", "en-au")
	req.Header.Set("Telstra-api-version", "3.x")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// requestWithQuery performs an HTTP request with query parameters.
func (c *Client) requestWithQuery(ctx context.Context, method, path string, query url.Values, result interface{}) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.request(ctx, method, path, nil, result)
}

// handleResponse processes the HTTP response.
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if result != nil && len(body) > 0 {
			if err := json.Unmarshal(body, result); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}
		}
		return nil
	}

	return c.parseError(resp.StatusCode, body)
}

// parseError parses an error response from the API.
func (c *Client) parseError(statusCode int, body []byte) error {
	var apiErrors APIErrors
	if err := json.Unmarshal(body, &apiErrors); err != nil {
		// If we can't parse the error, create a generic one
		return &APIError{
			Code:       "UNKNOWN_ERROR",
			Issue:      string(body),
			StatusCode: statusCode,
			Err:        sentinelForStatusCode(statusCode),
		}
	}

	sentinel := sentinelForStatusCode(statusCode)
	apiErrors.StatusCode = statusCode

	// Set the sentinel error on each error in the list
	for i := range apiErrors.Errors {
		apiErrors.Errors[i].StatusCode = statusCode
		apiErrors.Errors[i].Err = sentinel
	}

	return &apiErrors
}

// buildQuery creates URL query values from ListOptions.
func buildQuery(opts *ListOptions) url.Values {
	if opts == nil {
		return nil
	}

	query := url.Values{}
	if opts.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}
	if opts.Filter != "" {
		query.Set("filter", opts.Filter)
	}
	return query
}

// buildMessagesQuery creates URL query values from GetMessagesOptions.
func buildMessagesQuery(opts *GetMessagesOptions) url.Values {
	if opts == nil {
		return nil
	}

	query := url.Values{}
	if opts.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}
	if opts.Direction != "" {
		query.Set("direction", string(opts.Direction))
	}
	if opts.Status != "" {
		query.Set("status", string(opts.Status))
	}
	if opts.Filter != "" {
		query.Set("filter", opts.Filter)
	}
	if opts.Reverse {
		query.Set("reverse", "true")
	}
	if opts.StartTime != nil {
		query.Set("startTime", opts.StartTime.Format(time.RFC3339))
	}
	if opts.EndTime != nil {
		query.Set("endTime", opts.EndTime.Format(time.RFC3339))
	}
	return query
}
