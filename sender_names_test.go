package gopholler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestGetSenderNames(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sender-names" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"MySender", "CompanyName"})
	})
	defer server.Close()

	names, err := client.GetSenderNames(context.Background())
	if err != nil {
		t.Fatalf("GetSenderNames() error = %v", err)
	}
	if len(names) != 2 {
		t.Errorf("got %d sender names, want 2", len(names))
	}
	if names[0] != "MySender" {
		t.Errorf("first name = %q, want %q", names[0], "MySender")
	}
}

func TestGetSenderNames_Empty(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{})
	})
	defer server.Close()

	names, err := client.GetSenderNames(context.Background())
	if err != nil {
		t.Fatalf("GetSenderNames() error = %v", err)
	}
	if len(names) != 0 {
		t.Errorf("got %d sender names, want 0", len(names))
	}
}

func TestGetSenderNames_FreeTrial(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(APIErrors{
			Errors: []APIError{
				{
					Code:            "FREE_TRIAL_SENDERNAME_ERR",
					Issue:           "SenderNames are allowed only for paid accounts",
					SuggestedAction: "Move to a paid account to use this feature",
				},
			},
		})
	})
	defer server.Close()

	_, err := client.GetSenderNames(context.Background())
	if err == nil {
		t.Fatal("expected error for free trial account")
	}
	if !errors.Is(err, ErrPaymentRequired) {
		t.Errorf("expected ErrPaymentRequired, got %v", err)
	}
}
