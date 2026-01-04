package gopholler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetReports(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ReportsResponse{
			Reports: []Report{
				{ReportId: "report-1", ReportStatus: ReportStatusCompleted, ReportType: "messages"},
				{ReportId: "report-2", ReportStatus: ReportStatusQueued, ReportType: "messages"},
			},
			Paging: &Paging{TotalCount: 2},
		})
	})
	defer server.Close()

	resp, err := client.GetReports(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetReports() error = %v", err)
	}
	if len(resp.Reports) != 2 {
		t.Errorf("got %d reports, want 2", len(resp.Reports))
	}
	if resp.Reports[0].ReportStatus != ReportStatusCompleted {
		t.Errorf("first report status = %q, want %q", resp.Reports[0].ReportStatus, ReportStatusCompleted)
	}
}

func TestGetReports_WithOptions(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "20" {
			t.Errorf("limit = %q, want %q", r.URL.Query().Get("limit"), "20")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ReportsResponse{Reports: []Report{}})
	})
	defer server.Close()

	_, err := client.GetReports(context.Background(), &ListOptions{Limit: 20})
	if err != nil {
		t.Fatalf("GetReports() error = %v", err)
	}
}

func TestGetReport(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/reports/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Report{
			ReportId:     "report-123",
			ReportStatus: ReportStatusCompleted,
			ReportUrl:    "https://example.com/download/report-123.csv",
		})
	})
	defer server.Close()

	report, err := client.GetReport(context.Background(), "report-123")
	if err != nil {
		t.Fatalf("GetReport() error = %v", err)
	}
	if report.ReportId != "report-123" {
		t.Errorf("ReportId = %q, want %q", report.ReportId, "report-123")
	}
	if report.ReportUrl == "" {
		t.Error("expected ReportUrl to be set")
	}
}

func TestCreateMessagesReport(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reports/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req CreateMessagesReportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.StartDate != "2024-01-01" {
			t.Errorf("StartDate = %q, want %q", req.StartDate, "2024-01-01")
		}
		if req.EndDate != "2024-01-31" {
			t.Errorf("EndDate = %q, want %q", req.EndDate, "2024-01-31")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(ReportCreated{
			ReportId:          "new-report-123",
			ReportStatus:      ReportStatusQueued,
			ReportCallbackUrl: req.ReportCallbackUrl,
		})
	})
	defer server.Close()

	report, err := client.CreateMessagesReport(context.Background(), CreateMessagesReportRequest{
		StartDate:         "2024-01-01",
		EndDate:           "2024-01-31",
		ReportCallbackUrl: "https://example.com/report-callback",
	})

	if err != nil {
		t.Fatalf("CreateMessagesReport() error = %v", err)
	}
	if report.ReportId != "new-report-123" {
		t.Errorf("ReportId = %q, want %q", report.ReportId, "new-report-123")
	}
	if report.ReportStatus != ReportStatusQueued {
		t.Errorf("ReportStatus = %q, want %q", report.ReportStatus, ReportStatusQueued)
	}
}
