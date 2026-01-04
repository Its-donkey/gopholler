package gopholler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// testServer creates a mock server for testing API endpoints.
// It handles OAuth token requests automatically and routes other requests to the handler.
func testServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()

	mux := http.NewServeMux()

	// OAuth endpoint
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OAuth2Token{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   "3600",
		})
	})

	// API endpoints
	mux.HandleFunc("/", handler)

	server := httptest.NewServer(mux)

	client := NewClient("test-id", "test-secret",
		WithBaseURL(server.URL),
		WithOAuthURL(server.URL),
	)

	return server, client
}

func TestSendMessage(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}

		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.From != "0401234567" {
			t.Errorf("From = %q, want %q", req.From, "0401234567")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(MessageSent{
			MessageId: "msg-123",
			Status:    MessageStatusQueued,
			To:        "0412345678",
			From:      "0401234567",
		})
	})
	defer server.Close()

	msg, err := client.SendMessage(context.Background(), SendMessageRequest{
		To:             "0412345678",
		From:           "0401234567",
		MessageContent: "Hello!",
	})

	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	if msg.MessageId != "msg-123" {
		t.Errorf("MessageId = %v, want %q", msg.MessageId, "msg-123")
	}
	if msg.Status != MessageStatusQueued {
		t.Errorf("Status = %q, want %q", msg.Status, MessageStatusQueued)
	}
}

func TestSendMessage_Error(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIErrors{
			Errors: []APIError{
				{
					Code:            "FIELD_INVALID",
					Issue:           "The 'to' field is invalid",
					SuggestedAction: "Use a valid phone number",
					Field:           "to",
				},
			},
		})
	})
	defer server.Close()

	_, err := client.SendMessage(context.Background(), SendMessageRequest{
		To:   "invalid",
		From: "0401234567",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}

	var apiErrs *APIErrors
	if errors.As(err, &apiErrs) {
		if apiErrs.First().Field != "to" {
			t.Errorf("Field = %q, want %q", apiErrs.First().Field, "to")
		}
	} else {
		t.Error("expected error to be *APIErrors")
	}
}

func TestGetMessages(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		// Check query params
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("limit = %q, want %q", r.URL.Query().Get("limit"), "10")
		}
		if r.URL.Query().Get("direction") != "outgoing" {
			t.Errorf("direction = %q, want %q", r.URL.Query().Get("direction"), "outgoing")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessagesResponse{
			Messages: []MessageGet{
				{MessageId: "msg-1", Status: MessageStatusDelivered},
				{MessageId: "msg-2", Status: MessageStatusSent},
			},
			Paging: &Paging{TotalCount: 2},
		})
	})
	defer server.Close()

	resp, err := client.GetMessages(context.Background(), &GetMessagesOptions{
		Limit:     10,
		Direction: MessageDirectionOutgoing,
	})

	if err != nil {
		t.Fatalf("GetMessages() error = %v", err)
	}
	if len(resp.Messages) != 2 {
		t.Errorf("got %d messages, want 2", len(resp.Messages))
	}
}

func TestGetMessage(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/messages/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageGet{
			MessageId:      "msg-123",
			Status:         MessageStatusDelivered,
			MessageContent: "Hello!",
		})
	})
	defer server.Close()

	msg, err := client.GetMessage(context.Background(), "msg-123")
	if err != nil {
		t.Fatalf("GetMessage() error = %v", err)
	}
	if msg.MessageId != "msg-123" {
		t.Errorf("MessageId = %q, want %q", msg.MessageId, "msg-123")
	}
}

func TestUpdateMessage(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req UpdateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.MessageContent != "Updated content" {
			t.Errorf("MessageContent = %q, want %q", req.MessageContent, "Updated content")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageUpdate{
			MessageId:      "msg-123",
			Status:         MessageStatusQueued,
			MessageContent: "Updated content",
		})
	})
	defer server.Close()

	msg, err := client.UpdateMessage(context.Background(), "msg-123", UpdateMessageRequest{
		To:             "0412345678",
		From:           "0401234567",
		MessageContent: "Updated content",
	})

	if err != nil {
		t.Fatalf("UpdateMessage() error = %v", err)
	}
	if msg.MessageContent != "Updated content" {
		t.Errorf("MessageContent = %q, want %q", msg.MessageContent, "Updated content")
	}
}

func TestDeleteMessage(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.DeleteMessage(context.Background(), "msg-123")
	if err != nil {
		t.Fatalf("DeleteMessage() error = %v", err)
	}
}

func TestUpdateMessageTags(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var body struct {
			Tags []string `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Tags) != 2 {
			t.Errorf("got %d tags, want 2", len(body.Tags))
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.UpdateMessageTags(context.Background(), "msg-123", []string{"promo", "urgent"})
	if err != nil {
		t.Fatalf("UpdateMessageTags() error = %v", err)
	}
}
