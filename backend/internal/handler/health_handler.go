package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/service"
)

type HealthService interface {
	Status() service.HealthStatus
}

type HealthHandler struct {
	healthService HealthService
}

func NewHealthHandler(healthService HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, h.healthService.Status())
}
