package gopholler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetFreeTrialNumbers(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/free-trial-numbers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FreeTrialNumbers{
			FreeTrialNumbers: []string{"0412345678", "0498765432"},
		})
	})
	defer server.Close()

	resp, err := client.GetFreeTrialNumbers(context.Background())
	if err != nil {
		t.Fatalf("GetFreeTrialNumbers() error = %v", err)
	}
	if len(resp.FreeTrialNumbers) != 2 {
		t.Errorf("got %d free trial numbers, want 2", len(resp.FreeTrialNumbers))
	}
	if resp.FreeTrialNumbers[0] != "0412345678" {
		t.Errorf("first number = %q, want %q", resp.FreeTrialNumbers[0], "0412345678")
	}
}

func TestAddFreeTrialNumbers(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/free-trial-numbers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req AddFreeTrialNumbersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if len(req.FreeTrialNumbers) != 2 {
			t.Errorf("got %d numbers, want 2", len(req.FreeTrialNumbers))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(FreeTrialNumbers{
			FreeTrialNumbers: append([]string{"0400000000"}, req.FreeTrialNumbers...),
		})
	})
	defer server.Close()

	resp, err := client.AddFreeTrialNumbers(context.Background(), []string{"0411111111", "0422222222"})
	if err != nil {
		t.Fatalf("AddFreeTrialNumbers() error = %v", err)
	}
	if len(resp.FreeTrialNumbers) != 3 {
		t.Errorf("got %d free trial numbers, want 3", len(resp.FreeTrialNumbers))
	}
}
