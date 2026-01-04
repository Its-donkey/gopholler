package gopholler

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      APIError
		expected string
	}{
		{
			name: "error with field",
			err: APIError{
				Code:            "FIELD_INVALID",
				Issue:           "The field is invalid",
				SuggestedAction: "Check the field format",
				Field:           "to",
			},
			expected: "FIELD_INVALID: The field is invalid (field: to, suggested: Check the field format)",
		},
		{
			name: "error without field",
			err: APIError{
				Code:            "TOKEN_EXPIRED",
				Issue:           "The token has expired",
				SuggestedAction: "Refresh the token",
			},
			expected: "TOKEN_EXPIRED: The token has expired (suggested: Refresh the token)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("APIError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	err := &APIError{
		Code: "TOKEN_INVALID",
		Err:  ErrAuthentication,
	}

	if !errors.Is(err, ErrAuthentication) {
		t.Error("expected APIError to wrap ErrAuthentication")
	}

	if errors.Is(err, ErrNotFound) {
		t.Error("expected APIError to not match ErrNotFound")
	}
}

func TestAPIErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		errs     APIErrors
		expected string
	}{
		{
			name:     "no errors",
			errs:     APIErrors{},
			expected: "unknown API error",
		},
		{
			name: "single error",
			errs: APIErrors{
				Errors: []APIError{
					{Code: "ERR1", Issue: "Issue 1", SuggestedAction: "Fix it"},
				},
			},
			expected: "ERR1: Issue 1 (suggested: Fix it)",
		},
		{
			name: "multiple errors",
			errs: APIErrors{
				Errors: []APIError{
					{Code: "ERR1", Issue: "Issue 1", SuggestedAction: "Fix it"},
					{Code: "ERR2", Issue: "Issue 2", SuggestedAction: "Fix it too"},
				},
			},
			expected: "ERR1: Issue 1 (suggested: Fix it) (and 1 more errors)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errs.Error(); got != tt.expected {
				t.Errorf("APIErrors.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAPIErrors_Unwrap(t *testing.T) {
	errs := &APIErrors{
		Errors: []APIError{
			{Code: "ERR1", Err: ErrBadRequest},
		},
	}

	if !errors.Is(errs, ErrBadRequest) {
		t.Error("expected APIErrors to wrap ErrBadRequest")
	}
}

func TestAPIErrors_First(t *testing.T) {
	t.Run("with errors", func(t *testing.T) {
		errs := &APIErrors{
			Errors: []APIError{
				{Code: "ERR1"},
				{Code: "ERR2"},
			},
		}
		first := errs.First()
		if first == nil {
			t.Fatal("expected First() to return non-nil")
		}
		if first.Code != "ERR1" {
			t.Errorf("First().Code = %q, want %q", first.Code, "ERR1")
		}
	})

	t.Run("no errors", func(t *testing.T) {
		errs := &APIErrors{}
		if errs.First() != nil {
			t.Error("expected First() to return nil for empty errors")
		}
	})
}

func TestSentinelForStatusCode(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   error
	}{
		{400, ErrBadRequest},
		{401, ErrAuthentication},
		{402, ErrPaymentRequired},
		{403, ErrForbidden},
		{404, ErrNotFound},
		{429, ErrRateLimit},
		{500, ErrServer},
		{502, ErrServer},
		{503, ErrServer},
		{200, nil},
		{201, nil},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.statusCode)), func(t *testing.T) {
			got := sentinelForStatusCode(tt.statusCode)
			if got != tt.expected {
				t.Errorf("sentinelForStatusCode(%d) = %v, want %v", tt.statusCode, got, tt.expected)
			}
		})
	}
}
