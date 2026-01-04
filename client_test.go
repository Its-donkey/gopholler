package gopholler

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-id", "test-secret")

	if client.clientID != "test-id" {
		t.Errorf("clientID = %q, want %q", client.clientID, "test-id")
	}
	if client.clientSecret != "test-secret" {
		t.Errorf("clientSecret = %q, want %q", client.clientSecret, "test-secret")
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", client.baseURL, DefaultBaseURL)
	}
	if client.oauthURL != DefaultOAuthURL {
		t.Errorf("oauthURL = %q, want %q", client.oauthURL, DefaultOAuthURL)
	}
	if len(client.scopes) == 0 {
		t.Error("expected default scopes to be set")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	customHTTP := &http.Client{Timeout: 60 * time.Second}

	client := NewClient("test-id", "test-secret",
		WithBaseURL("https://custom.api.com/v3"),
		WithOAuthURL("https://custom.oauth.com"),
		WithHTTPClient(customHTTP),
		WithScopes([]string{"messages:read"}),
	)

	if client.baseURL != "https://custom.api.com/v3" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "https://custom.api.com/v3")
	}
	if client.oauthURL != "https://custom.oauth.com" {
		t.Errorf("oauthURL = %q, want %q", client.oauthURL, "https://custom.oauth.com")
	}
	if client.httpClient != customHTTP {
		t.Error("expected custom HTTP client to be set")
	}
	if len(client.scopes) != 1 || client.scopes[0] != "messages:read" {
		t.Errorf("scopes = %v, want [messages:read]", client.scopes)
	}
}

func TestWithBaseURL_TrailingSlash(t *testing.T) {
	client := NewClient("id", "secret", WithBaseURL("https://api.com/"))
	if client.baseURL != "https://api.com" {
		t.Errorf("baseURL = %q, want trailing slash removed", client.baseURL)
	}
}

func TestBuildQuery(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		query := buildQuery(nil)
		if query != nil {
			t.Error("expected nil query for nil options")
		}
	})

	t.Run("empty options", func(t *testing.T) {
		query := buildQuery(&ListOptions{})
		if len(query) != 0 {
			t.Errorf("expected empty query, got %v", query)
		}
	})

	t.Run("with values", func(t *testing.T) {
		query := buildQuery(&ListOptions{
			Limit:  10,
			Offset: 20,
			Filter: "tag:test",
		})

		if query.Get("limit") != "10" {
			t.Errorf("limit = %q, want %q", query.Get("limit"), "10")
		}
		if query.Get("offset") != "20" {
			t.Errorf("offset = %q, want %q", query.Get("offset"), "20")
		}
		if query.Get("filter") != "tag:test" {
			t.Errorf("filter = %q, want %q", query.Get("filter"), "tag:test")
		}
	})
}

func TestBuildMessagesQuery(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		query := buildMessagesQuery(nil)
		if query != nil {
			t.Error("expected nil query for nil options")
		}
	})

	t.Run("with all values", func(t *testing.T) {
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

		query := buildMessagesQuery(&GetMessagesOptions{
			Limit:     50,
			Offset:    100,
			Direction: MessageDirectionOutgoing,
			Status:    MessageStatusDelivered,
			Filter:    "campaign:promo",
			Reverse:   true,
			StartTime: &startTime,
			EndTime:   &endTime,
		})

		if query.Get("limit") != "50" {
			t.Errorf("limit = %q, want %q", query.Get("limit"), "50")
		}
		if query.Get("offset") != "100" {
			t.Errorf("offset = %q, want %q", query.Get("offset"), "100")
		}
		if query.Get("direction") != "outgoing" {
			t.Errorf("direction = %q, want %q", query.Get("direction"), "outgoing")
		}
		if query.Get("status") != "delivered" {
			t.Errorf("status = %q, want %q", query.Get("status"), "delivered")
		}
		if query.Get("filter") != "campaign:promo" {
			t.Errorf("filter = %q, want %q", query.Get("filter"), "campaign:promo")
		}
		if query.Get("reverse") != "true" {
			t.Errorf("reverse = %q, want %q", query.Get("reverse"), "true")
		}
		if query.Get("startTime") == "" {
			t.Error("expected startTime to be set")
		}
		if query.Get("endTime") == "" {
			t.Error("expected endTime to be set")
		}
	})
}
