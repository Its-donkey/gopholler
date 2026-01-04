package gopholler

import (
	"context"
	"net/http"
	"net/url"
)

// GetLogs retrieves API call logs for the account.
// Returns the last 50 calls by default, or filtered by date range/resource.
func (c *Client) GetLogs(ctx context.Context, opts *GetLogsOptions) (*LogsResponse, error) {
	query := buildLogsQuery(opts)
	var result LogsResponse
	err := c.requestWithQuery(ctx, http.MethodGet, "/logs", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// buildLogsQuery creates URL query values from GetLogsOptions.
func buildLogsQuery(opts *GetLogsOptions) url.Values {
	if opts == nil {
		return nil
	}

	query := url.Values{}
	if opts.StartDate != "" {
		query.Set("startDate", opts.StartDate)
	}
	if opts.EndDate != "" {
		query.Set("endDate", opts.EndDate)
	}
	if opts.Filter != "" {
		query.Set("filter", opts.Filter)
	}
	return query
}
