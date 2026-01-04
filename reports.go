package gopholler

import (
	"context"
	"fmt"
	"net/http"
)

// GetReports retrieves all reports for the account.
func (c *Client) GetReports(ctx context.Context, opts *ListOptions) (*ReportsResponse, error) {
	query := buildQuery(opts)
	var result ReportsResponse
	err := c.requestWithQuery(ctx, http.MethodGet, "/reports", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetReport retrieves a specific report by ID.
func (c *Client) GetReport(ctx context.Context, reportID string) (*Report, error) {
	path := fmt.Sprintf("/reports/%s", reportID)
	var result Report
	err := c.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateMessagesReport creates a new messages report.
func (c *Client) CreateMessagesReport(ctx context.Context, req CreateMessagesReportRequest) (*ReportCreated, error) {
	var result ReportCreated
	err := c.request(ctx, http.MethodPost, "/reports/messages", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
