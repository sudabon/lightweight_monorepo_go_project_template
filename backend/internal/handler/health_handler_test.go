package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/service"
)

func TestHealthHandlerGet(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://invalid:invalid@127.0.0.1:1/invalid")

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
