package gopholler

import (
	"context"
	"net/http"
)

// GetSenderNames retrieves all approved sender names for the account.
// This is a paid feature - free trial accounts will receive a 402 error.
func (c *Client) GetSenderNames(ctx context.Context) ([]string, error) {
	var result []string
	err := c.request(ctx, http.MethodGet, "/sender-names", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
