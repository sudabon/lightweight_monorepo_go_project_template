package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/config"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/db"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/handler"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/repository"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/router"
	"github.com/sudabon/lightweight_monorepo_go_project_template/backend/internal/service"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg := config.Load()
	pool, err := db.OpenPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer pool.Close()

	healthRepository := repository.NewHealthRepository(pool)
	healthService := service.NewHealthService()
	dbHealthService := service.NewDBHealthService(healthRepository)
	healthHandler := handler.NewHealthHandler(healthService, dbHealthService)
	corsAllowOrigins := strings.Split(cfg.CORSAllowOrigins, ",")
	for i := range corsAllowOrigins {
		corsAllowOrigins[i] = strings.TrimSpace(corsAllowOrigins[i])
	}
	e := router.New(healthHandler, corsAllowOrigins)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := e.Start(":" + cfg.AppPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		if err == nil {
			return
		}
		log.Fatalf("server failed: %v", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("server shutdown failed: %v", err)
		}
		if err := <-errCh; err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}
}
