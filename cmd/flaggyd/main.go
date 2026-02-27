package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/alexis/flaggy/internal/api"
	"github.com/alexis/flaggy/internal/config"
	"github.com/alexis/flaggy/internal/sse"
	"github.com/alexis/flaggy/internal/store"
	"github.com/alexis/flaggy/migrations"
)

var version = "dev"

func main() {
	_ = godotenv.Load() // .env is optional

	cfg := config.Load()

	slog.Info("starting flaggy", "version", version, "port", cfg.Port, "db", cfg.DBPath)

	db, err := store.NewSQLiteStore(cfg.DBPath, migrations.FS)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	broadcaster := sse.NewBroadcaster()
	defer broadcaster.Close()

	if cfg.MasterKey == "" {
		slog.Warn("FLAGGY_MASTER_KEY not set — auth disabled (dev mode)")
	}

	router := api.NewRouter(db, broadcaster, cfg.MasterKey)

	srv := &http.Server{
		Addr:        cfg.Port,
		Handler:     router,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
		// No WriteTimeout — SSE streams are long-lived connections
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server listening", "addr", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}

	slog.Info("server stopped")
}
