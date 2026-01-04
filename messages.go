package gopholler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// IsValidMobileNumber checks if a phone number is a valid mobile number.
// It validates against known mobile prefixes for countries worldwide.
//
// Valid formats:
//   - International: +[country code][mobile prefix][number] (e.g., +61412345678)
//   - Australian domestic: 04XXXXXXXX (10 digits starting with 04)
//
// Spaces and dashes are stripped before validation.
func IsValidMobileNumber(number string) bool {
	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	if len(cleaned) < 8 {
		return false
	}

	// Handle international format: +[country code][number]
	if strings.HasPrefix(cleaned, "+") {
		return isValidInternationalMobile(cleaned[1:])
	}

	// Handle domestic formats
	for countryCode, domesticPrefix := range domesticMobilePrefixes {
		if strings.HasPrefix(cleaned, domesticPrefix) {
			// Convert to international format and validate
			// e.g., Australian 0412345678 -> 614XXXXXXXX
			internationalNum := countryCode + cleaned[1:]
			return isValidInternationalMobile(internationalNum)
		}
	}

	return false
}

// isValidInternationalMobile validates an international mobile number (without the + prefix).
// It extracts the country code and checks if the number starts with a valid mobile prefix.
func isValidInternationalMobile(number string) bool {
	if !isAllDigits(number) {
		return false
	}

	// Try to match country codes (1-3 digits)
	for codeLen := 1; codeLen <= 3 && codeLen < len(number); codeLen++ {
		countryCode := number[:codeLen]
		remainder := number[codeLen:]

		prefixes, exists := mobilePrefixes[countryCode]
		if !exists {
			continue
		}

		// Check if the remainder starts with any valid mobile prefix
		for _, prefix := range prefixes {
			if strings.HasPrefix(remainder, prefix) {
				return true
			}
		}
	}

	return false
}

// IsValidAustralianMobile checks if a phone number is a valid Australian mobile number.
// Valid formats:
//   - Domestic: 04XX XXX XXX (10 digits starting with 04)
//   - International: +614XXXXXXXX (12 characters starting with +614)
//
// Spaces and dashes are stripped before validation.
// Deprecated: Use IsValidMobileNumber instead, which supports all countries.
func IsValidAustralianMobile(number string) bool {
	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Check international format: +614XXXXXXXX (12 chars)
	if strings.HasPrefix(cleaned, "+614") && len(cleaned) == 12 {
		return isAllDigits(cleaned[1:]) // Skip the + sign
	}

	// Check domestic format: 04XXXXXXXX (10 digits)
	if strings.HasPrefix(cleaned, "04") && len(cleaned) == 10 {
		return isAllDigits(cleaned)
	}

	return false
}

// isAllDigits returns true if s contains only digit characters.
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// validateMobileNumbers validates that all recipient numbers are valid mobile numbers.
func validateMobileNumbers(to any) error {
	switch v := to.(type) {
	case string:
		if !IsValidMobileNumber(v) {
			return fmt.Errorf("%w: %s", ErrInvalidMobileNumber, v)
		}
	case []string:
		for _, num := range v {
			if !IsValidMobileNumber(num) {
				return fmt.Errorf("%w: %s", ErrInvalidMobileNumber, num)
			}
		}
	case []any:
		for _, item := range v {
			if num, ok := item.(string); ok {
				if !IsValidMobileNumber(num) {
					return fmt.Errorf("%w: %s", ErrInvalidMobileNumber, num)
				}
			}
		}
	}
	return nil
}

// SendMessage sends an SMS or MMS message.
func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*MessageSent, error) {
	if err := validateMobileNumbers(req.To); err != nil {
		return nil, err
	}

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
