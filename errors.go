package gopholler

import (
	"errors"
	"fmt"
)

// Sentinel errors for use with errors.Is()
var (
	ErrAuthentication      = errors.New("authentication failed")
	ErrRateLimit           = errors.New("rate limit exceeded")
	ErrNotFound            = errors.New("resource not found")
	ErrBadRequest          = errors.New("bad request")
	ErrForbidden           = errors.New("forbidden")
	ErrPaymentRequired     = errors.New("payment required")
	ErrServer              = errors.New("server error")
	ErrInvalidMobileNumber = errors.New("invalid mobile number")
)

// APIError represents an error returned by the Telstra Messaging API.
type APIError struct {
	Code            string `json:"code"`
	Issue           string `json:"issue"`
	SuggestedAction string `json:"suggested_action"`
	Field           string `json:"field,omitempty"`
	Value           string `json:"value,omitempty"`
	Location        string `json:"location,omitempty"`
	Link            string `json:"link,omitempty"`
	StatusCode      int    `json:"-"`
	Err             error  `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s, suggested: %s)", e.Code, e.Issue, e.Field, e.SuggestedAction)
	}
	return fmt.Sprintf("%s: %s (suggested: %s)", e.Code, e.Issue, e.SuggestedAction)
}

// Unwrap returns the underlying sentinel error for use with errors.Is().
func (e *APIError) Unwrap() error {
	return e.Err
}

// APIErrors represents a collection of API errors.
type APIErrors struct {
	Errors     []APIError `json:"errors"`
	StatusCode int        `json:"-"`
}

// Error implements the error interface.
func (e *APIErrors) Error() string {
	if len(e.Errors) == 0 {
		return "unknown API error"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", e.Errors[0].Error(), len(e.Errors)-1)
}

// Unwrap returns the first error's underlying sentinel error for use with errors.Is().
func (e *APIErrors) Unwrap() error {
	if len(e.Errors) > 0 {
		return e.Errors[0].Err
	}
	return nil
}

// First returns the first API error, or nil if there are no errors.
func (e *APIErrors) First() *APIError {
	if len(e.Errors) > 0 {
		return &e.Errors[0]
	}
	return nil
}

// sentinelForStatusCode returns the appropriate sentinel error for an HTTP status code.
func sentinelForStatusCode(statusCode int) error {
	switch statusCode {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrAuthentication
	case 402:
		return ErrPaymentRequired
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 429:
		return ErrRateLimit
	default:
		if statusCode >= 500 {
			return ErrServer
		}
		return nil
	}
}
