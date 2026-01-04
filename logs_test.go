package gopholler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGetLogs(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/logs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LogsResponse{
			Logs: []Log{
				{Timestamp: &timestamp, Value: "POST /messages", StatusCode: 201},
				{Timestamp: &timestamp, Value: "GET /messages", StatusCode: 200},
			},
		})
	})
	defer server.Close()

	resp, err := client.GetLogs(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
	if len(resp.Logs) != 2 {
		t.Errorf("got %d logs, want 2", len(resp.Logs))
	}
	if resp.Logs[0].StatusCode != 201 {
		t.Errorf("first log status = %d, want 201", resp.Logs[0].StatusCode)
	}
}

func TestGetLogs_WithOptions(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("startDate") != "2024-01-01T00:00:00Z" {
			t.Errorf("startDate = %q, want %q", r.URL.Query().Get("startDate"), "2024-01-01T00:00:00Z")
		}
		if r.URL.Query().Get("endDate") != "2024-01-31T23:59:59Z" {
			t.Errorf("endDate = %q, want %q", r.URL.Query().Get("endDate"), "2024-01-31T23:59:59Z")
		}
		if r.URL.Query().Get("filter") != "POST /messages" {
			t.Errorf("filter = %q, want %q", r.URL.Query().Get("filter"), "POST /messages")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LogsResponse{Logs: []Log{}})
	})
	defer server.Close()

	_, err := client.GetLogs(context.Background(), &GetLogsOptions{
		StartDate: "2024-01-01T00:00:00Z",
		EndDate:   "2024-01-31T23:59:59Z",
		Filter:    "POST /messages",
	})
	if err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
}

func TestGetLogs_Empty(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LogsResponse{})
	})
	defer server.Close()

	resp, err := client.GetLogs(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
	if len(resp.Logs) != 0 {
		t.Errorf("got %d logs, want 0", len(resp.Logs))
	}
}
