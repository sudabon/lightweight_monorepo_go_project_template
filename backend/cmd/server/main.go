package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/config"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/router"
)

func main() {
	cfg := config.Load()
	e := router.New()

	if err := e.Start(":" + cfg.AppPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
