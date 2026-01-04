package gopholler

import (
	"context"
	"fmt"
	"net/http"
)

// GetVirtualNumbers retrieves all virtual numbers assigned to the account.
func (c *Client) GetVirtualNumbers(ctx context.Context, opts *ListOptions) (*VirtualNumbersResponse, error) {
	query := buildQuery(opts)
	var result VirtualNumbersResponse
	err := c.requestWithQuery(ctx, http.MethodGet, "/virtual-numbers", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AssignVirtualNumber assigns a new virtual number to the account.
func (c *Client) AssignVirtualNumber(ctx context.Context, req AssignVirtualNumberRequest) (*VirtualNumber, error) {
	var result VirtualNumber
	err := c.request(ctx, http.MethodPost, "/virtual-numbers", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVirtualNumber retrieves a specific virtual number.
func (c *Client) GetVirtualNumber(ctx context.Context, virtualNumber string) (*VirtualNumber, error) {
	path := fmt.Sprintf("/virtual-numbers/%s", virtualNumber)
	var result VirtualNumber
	err := c.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateVirtualNumber updates a virtual number's configuration.
func (c *Client) UpdateVirtualNumber(ctx context.Context, virtualNumber string, req UpdateVirtualNumberRequest) (*VirtualNumber, error) {
	path := fmt.Sprintf("/virtual-numbers/%s", virtualNumber)
	var result VirtualNumber
	err := c.request(ctx, http.MethodPut, path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteVirtualNumber removes a virtual number from the account.
func (c *Client) DeleteVirtualNumber(ctx context.Context, virtualNumber string) error {
	path := fmt.Sprintf("/virtual-numbers/%s", virtualNumber)
	return c.request(ctx, http.MethodDelete, path, nil, nil)
}

// GetVirtualNumberOptouts retrieves opt-outs for a virtual number.
func (c *Client) GetVirtualNumberOptouts(ctx context.Context, virtualNumber string, opts *ListOptions) (*OptoutsResponse, error) {
	path := fmt.Sprintf("/virtual-numbers/%s/optouts", virtualNumber)
	query := buildQuery(opts)
	var result OptoutsResponse
	err := c.requestWithQuery(ctx, http.MethodGet, path, query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
