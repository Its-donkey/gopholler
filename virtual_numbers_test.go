package gopholler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetVirtualNumbers(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/virtual-numbers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VirtualNumbersResponse{
			VirtualNumbers: []VirtualNumber{
				{VirtualNumber: "0401234567"},
				{VirtualNumber: "0412345678"},
			},
			Paging: &Paging{TotalCount: 2},
		})
	})
	defer server.Close()

	resp, err := client.GetVirtualNumbers(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetVirtualNumbers() error = %v", err)
	}
	if len(resp.VirtualNumbers) != 2 {
		t.Errorf("got %d virtual numbers, want 2", len(resp.VirtualNumbers))
	}
}

func TestGetVirtualNumbers_WithOptions(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("limit = %q, want %q", r.URL.Query().Get("limit"), "5")
		}
		if r.URL.Query().Get("offset") != "10" {
			t.Errorf("offset = %q, want %q", r.URL.Query().Get("offset"), "10")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VirtualNumbersResponse{
			VirtualNumbers: []VirtualNumber{},
		})
	})
	defer server.Close()

	_, err := client.GetVirtualNumbers(context.Background(), &ListOptions{
		Limit:  5,
		Offset: 10,
	})
	if err != nil {
		t.Fatalf("GetVirtualNumbers() error = %v", err)
	}
}

func TestAssignVirtualNumber(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/virtual-numbers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req AssignVirtualNumberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.ReplyCallbackUrl != "https://example.com/callback" {
			t.Errorf("ReplyCallbackUrl = %q, want %q", req.ReplyCallbackUrl, "https://example.com/callback")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(VirtualNumber{
			VirtualNumber:    "0401234567",
			ReplyCallbackUrl: "https://example.com/callback",
		})
	})
	defer server.Close()

	vn, err := client.AssignVirtualNumber(context.Background(), AssignVirtualNumberRequest{
		ReplyCallbackUrl: "https://example.com/callback",
		Tags:             []string{"test"},
	})

	if err != nil {
		t.Fatalf("AssignVirtualNumber() error = %v", err)
	}
	if vn.VirtualNumber != "0401234567" {
		t.Errorf("VirtualNumber = %q, want %q", vn.VirtualNumber, "0401234567")
	}
}

func TestGetVirtualNumber(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/virtual-numbers/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VirtualNumber{
			VirtualNumber: "0401234567",
		})
	})
	defer server.Close()

	vn, err := client.GetVirtualNumber(context.Background(), "0401234567")
	if err != nil {
		t.Fatalf("GetVirtualNumber() error = %v", err)
	}
	if vn.VirtualNumber != "0401234567" {
		t.Errorf("VirtualNumber = %q, want %q", vn.VirtualNumber, "0401234567")
	}
}

func TestUpdateVirtualNumber(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req UpdateVirtualNumberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VirtualNumber{
			VirtualNumber:    "0401234567",
			ReplyCallbackUrl: req.ReplyCallbackUrl,
		})
	})
	defer server.Close()

	vn, err := client.UpdateVirtualNumber(context.Background(), "0401234567", UpdateVirtualNumberRequest{
		ReplyCallbackUrl: "https://new-callback.com",
	})

	if err != nil {
		t.Fatalf("UpdateVirtualNumber() error = %v", err)
	}
	if vn.ReplyCallbackUrl != "https://new-callback.com" {
		t.Errorf("ReplyCallbackUrl = %q, want %q", vn.ReplyCallbackUrl, "https://new-callback.com")
	}
}

func TestDeleteVirtualNumber(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteVirtualNumber(context.Background(), "0401234567")
	if err != nil {
		t.Fatalf("DeleteVirtualNumber() error = %v", err)
	}
}

func TestGetVirtualNumberOptouts(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/optouts") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OptoutsResponse{
			RecipientOptouts: []RecipientOptout{
				{OptoutNumber: "0412345678"},
			},
			Paging: &Paging{TotalCount: 1},
		})
	})
	defer server.Close()

	resp, err := client.GetVirtualNumberOptouts(context.Background(), "0401234567", nil)
	if err != nil {
		t.Fatalf("GetVirtualNumberOptouts() error = %v", err)
	}
	if len(resp.RecipientOptouts) != 1 {
		t.Errorf("got %d optouts, want 1", len(resp.RecipientOptouts))
	}
}
