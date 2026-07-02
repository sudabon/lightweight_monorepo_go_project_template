package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/service"
)

type fakeDBHealthService struct {
	status service.DBHealthStatus
}

func (f fakeDBHealthService) Check(ctx context.Context) service.DBHealthStatus {
	return f.status
}

func TestHealthHandlerGet_DB非依存_OKを返す(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHealthHandler(service.NewHealthService())

	if err := h.Get(c); err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("response JSON decode failed: %v", err)
	}

	if len(got) != 1 || got["status"] != "ok" {
		t.Fatalf("response body = %#v, want %#v", got, map[string]string{"status": "ok"})
	}
}

func TestHealthHandlerGetDB_DBステータス_HTTPレスポンスを返す(t *testing.T) {
	tests := []struct {
		name     string
		status   service.DBHealthStatus
		wantCode int
		wantBody map[string]string
	}{
		{
			name:     "db ok",
			status:   service.DBHealthStatus{Status: "ok"},
			wantCode: http.StatusOK,
			wantBody: map[string]string{"status": "ok"},
		},
		{
			name:     "db unavailable",
			status:   service.DBHealthStatus{Status: "unavailable"},
			wantCode: http.StatusServiceUnavailable,
			wantBody: map[string]string{"status": "unavailable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/health/db", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := NewHealthHandler(service.NewHealthService(), fakeDBHealthService{status: tt.status})

			if err := h.GetDB(c); err != nil {
				t.Fatalf("GetDB() returned error: %v", err)
			}

			if rec.Code != tt.wantCode {
				t.Fatalf("status code = %d, want %d", rec.Code, tt.wantCode)
			}

			var got map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("response JSON decode failed: %v", err)
			}

			if len(got) != 1 || got["status"] != tt.wantBody["status"] {
				t.Fatalf("response body = %#v, want %#v", got, tt.wantBody)
			}
		})
	}
}
