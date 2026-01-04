package gopholler

import (
	"context"
	"net/http"
)

// GetFreeTrialNumbers retrieves all free trial numbers for the account.
func (c *Client) GetFreeTrialNumbers(ctx context.Context) (*FreeTrialNumbers, error) {
	var result FreeTrialNumbers
	err := c.request(ctx, http.MethodGet, "/free-trial-numbers", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddFreeTrialNumbers adds numbers to the free trial numbers list.
func (c *Client) AddFreeTrialNumbers(ctx context.Context, numbers []string) (*FreeTrialNumbers, error) {
	req := AddFreeTrialNumbersRequest{
		FreeTrialNumbers: numbers,
	}
	var result FreeTrialNumbers
	err := c.request(ctx, http.MethodPost, "/free-trial-numbers", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
