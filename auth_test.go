package gopholler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestAuthenticate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("client_id") != "test-id" {
			t.Errorf("client_id = %q, want %q", r.FormValue("client_id"), "test-id")
		}
		if r.FormValue("client_secret") != "test-secret" {
			t.Errorf("client_secret = %q, want %q", r.FormValue("client_secret"), "test-secret")
		}
		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("grant_type = %q, want %q", r.FormValue("grant_type"), "client_credentials")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OAuth2Token{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   "3600",
		})
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithOAuthURL(server.URL))

	token, err := client.authenticate(context.Background())
	if err != nil {
		t.Fatalf("authenticate() error = %v", err)
	}
	if token.AccessToken != "test-token" {
		t.Errorf("AccessToken = %q, want %q", token.AccessToken, "test-token")
	}
}

func TestAuthenticate_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
	}))
	defer server.Close()

	client := NewClient("bad-id", "bad-secret", WithOAuthURL(server.URL))

	_, err := client.authenticate(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid credentials")
	}
	if !errors.Is(err, ErrAuthentication) {
		t.Errorf("expected ErrAuthentication, got %v", err)
	}
}

func TestGetValidToken_CachesToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OAuth2Token{
			AccessToken: "cached-token",
			TokenType:   "Bearer",
			ExpiresIn:   "3600",
		})
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithOAuthURL(server.URL))

	// First call should authenticate
	token1, err := client.getValidToken(context.Background())
	if err != nil {
		t.Fatalf("first getValidToken() error = %v", err)
	}

	// Second call should use cache
	token2, err := client.getValidToken(context.Background())
	if err != nil {
		t.Fatalf("second getValidToken() error = %v", err)
	}

	if token1.AccessToken != token2.AccessToken {
		t.Error("expected same token from cache")
	}

	if callCount != 1 {
		t.Errorf("expected 1 auth call, got %d", callCount)
	}
}

func TestGetValidToken_RefreshesExpiredToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OAuth2Token{
			AccessToken: "token-" + string(rune('0'+callCount)),
			TokenType:   "Bearer",
			ExpiresIn:   "1", // Expires in 1 second
		})
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithOAuthURL(server.URL))

	// First call
	_, err := client.getValidToken(context.Background())
	if err != nil {
		t.Fatalf("first getValidToken() error = %v", err)
	}

	// Force expiry
	client.tokenMu.Lock()
	client.tokenExpiry = time.Now().Add(-time.Hour)
	client.tokenMu.Unlock()

	// Second call should refresh
	_, err = client.getValidToken(context.Background())
	if err != nil {
		t.Fatalf("second getValidToken() error = %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 auth calls after expiry, got %d", callCount)
	}
}

func TestGetValidToken_ThreadSafe(t *testing.T) {
	callCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		callCount++
		mu.Unlock()

		// Simulate slow auth
		time.Sleep(50 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OAuth2Token{
			AccessToken: "concurrent-token",
			TokenType:   "Bearer",
			ExpiresIn:   "3600",
		})
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithOAuthURL(server.URL))

	// Launch multiple concurrent requests
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.getValidToken(context.Background())
			if err != nil {
				t.Errorf("getValidToken() error = %v", err)
			}
		}()
	}
	wg.Wait()

	// Should only have made 1-2 auth calls due to locking
	mu.Lock()
	defer mu.Unlock()
	if callCount > 2 {
		t.Errorf("expected at most 2 auth calls with concurrency, got %d", callCount)
	}
}

func TestInvalidateToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OAuth2Token{
			AccessToken: "token",
			TokenType:   "Bearer",
			ExpiresIn:   "3600",
		})
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithOAuthURL(server.URL))

	// Get initial token
	_, err := client.getValidToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Invalidate
	client.InvalidateToken()

	// Should fetch new token
	_, err = client.getValidToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 auth calls after invalidate, got %d", callCount)
	}
}
