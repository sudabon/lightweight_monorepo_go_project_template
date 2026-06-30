package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/handler"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/service"
)

func New() *echo.Echo {
	e := echo.New()

	healthService := service.NewHealthService()
	healthHandler := handler.NewHealthHandler(healthService)
	e.GET("/health", healthHandler.Get)

	return e
}
