package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/handler"
)

// New はEchoインスタンスを構成し、ルートを登録します。
func New(healthHandler *handler.HealthHandler, corsAllowOrigins []string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsAllowOrigins,
	}))

	e.GET("/health", healthHandler.Get)
	e.GET("/health/db", healthHandler.GetDB)

	return e
}
