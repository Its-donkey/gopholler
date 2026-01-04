package gopholler

import (
	"context"
	"net/http"
)

// HealthCheck checks if the API is available.
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.request(ctx, http.MethodGet, "/health-check", nil, nil)
}
