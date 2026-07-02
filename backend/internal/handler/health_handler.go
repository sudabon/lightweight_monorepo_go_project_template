package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/service"
)

// HealthService は通常のヘルスチェックに必要な振る舞いを表します。
type HealthService interface {
	Status() service.HealthStatus
}

// DBHealthService はDBヘルスチェックに必要な振る舞いを表します。
type DBHealthService interface {
	Check(ctx context.Context) service.DBHealthStatus
}

// HealthHandler はヘルスチェックHTTPリクエストを扱います。
type HealthHandler struct {
	healthService   HealthService
	dbHealthService DBHealthService
}

// NewHealthHandler はHealthHandlerを生成します。
func NewHealthHandler(healthService HealthService, dbHealthServices ...DBHealthService) *HealthHandler {
	var dbHealthService DBHealthService
	if len(dbHealthServices) > 0 {
		dbHealthService = dbHealthServices[0]
	}

	return &HealthHandler{
		healthService:   healthService,
		dbHealthService: dbHealthService,
	}
}

// Get はDB非依存のヘルスチェックを返します。
func (h *HealthHandler) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, h.healthService.Status())
}

// GetDB はDB疎通確認のヘルスチェックを返します。
func (h *HealthHandler) GetDB(c echo.Context) error {
	if h.dbHealthService == nil {
		return c.JSON(http.StatusServiceUnavailable, service.DBHealthStatus{Status: "unavailable"})
	}

	status := h.dbHealthService.Check(c.Request().Context())
	if status.Status != "ok" {
		return c.JSON(http.StatusServiceUnavailable, status)
	}

	return c.JSON(http.StatusOK, status)
}
