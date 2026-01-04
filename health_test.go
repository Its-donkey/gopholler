package gopholler

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestHealthCheck_Success(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health-check" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}
}

func TestHealthCheck_ServiceUnavailable(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"errors":[{"code":"SERVICE_UNAVAILABLE","issue":"Service temporarily unavailable","suggested_action":"Try again later"}]}`))
	})
	defer server.Close()

	err := client.HealthCheck(context.Background())
	if err == nil {
		t.Fatal("expected error for service unavailable")
	}
	if !errors.Is(err, ErrServer) {
		t.Errorf("expected ErrServer, got %v", err)
	}
}
