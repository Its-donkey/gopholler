package gopholler

import (
	"context"
	"fmt"
	"net/http"
)

// SendMessage sends an SMS or MMS message.
func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*MessageSent, error) {
	var result MessageSent
	err := c.request(ctx, http.MethodPost, "/messages", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMessages retrieves all sent/received messages.
func (c *Client) GetMessages(ctx context.Context, opts *GetMessagesOptions) (*MessagesResponse, error) {
	query := buildMessagesQuery(opts)
	var result MessagesResponse
	err := c.requestWithQuery(ctx, http.MethodGet, "/messages", query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMessage retrieves a specific message by ID.
func (c *Client) GetMessage(ctx context.Context, messageID string) (*MessageGet, error) {
	path := fmt.Sprintf("/messages/%s", messageID)
	var result MessageGet
	err := c.request(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateMessage updates a scheduled message that hasn't been sent yet.
func (c *Client) UpdateMessage(ctx context.Context, messageID string, req UpdateMessageRequest) (*MessageUpdate, error) {
	path := fmt.Sprintf("/messages/%s", messageID)
	var result MessageUpdate
	err := c.request(ctx, http.MethodPut, path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMessage deletes a scheduled message that hasn't been sent yet.
func (c *Client) DeleteMessage(ctx context.Context, messageID string) error {
	path := fmt.Sprintf("/messages/%s", messageID)
	return c.request(ctx, http.MethodDelete, path, nil, nil)
}

// UpdateMessageTags updates the tags on a message.
func (c *Client) UpdateMessageTags(ctx context.Context, messageID string, tags []string) error {
	path := fmt.Sprintf("/messages/%s", messageID)
	body := struct {
		Tags []string `json:"tags"`
	}{Tags: tags}
	return c.request(ctx, http.MethodPatch, path, body, nil)
}
