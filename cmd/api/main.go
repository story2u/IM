package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/httpapi"
	"im-go/internal/integrationhub/service"
	"im-go/internal/integrationhub/store"
)

type serverConfig struct {
	Addr        string
	DatabaseURL string
	AutoMigrate bool
	SeedData    bool
}

func main() {
	cfg := loadConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	handler, cleanup, err := buildHandler(context.Background(), cfg, time.Now().UTC(), logger)
	if err != nil {
		logger.Error("initialize integration API failed", "error", err)
		os.Exit(1)
	}
	defer cleanup()
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting integration API", "addr", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("shutdown complete")
}

func buildHandler(ctx context.Context, cfg serverConfig, now time.Time, logger *slog.Logger) (http.Handler, func(), error) {
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		logger.Warn("IM_DATABASE_URL is not configured; using in-memory integration repository for local development")
		repo := store.NewMemory(domain.SeedSnapshot(now))
		svc := service.New(repo)
		return httpapi.New(svc), func() {}, nil
	}
	repo, err := store.OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { repo.Close() }
	if cfg.AutoMigrate {
		if err := repo.Migrate(ctx); err != nil {
			cleanup()
			return nil, nil, err
		}
	}
	if cfg.SeedData {
		if err := repo.SeedIfEmpty(ctx, domain.SeedSnapshot(now)); err != nil {
			cleanup()
			return nil, nil, err
		}
	}
	svc := service.New(repo)
	return httpapi.New(svc), cleanup, nil
}

func loadConfig() serverConfig {
	addr := firstEnv("IM_API_ADDR", "GO_BACKEND_ADDR", "ADDR")
	if addr == "" {
		addr = ":9000"
	}
	return serverConfig{
		Addr:        addr,
		DatabaseURL: firstEnv("IM_DATABASE_URL", "INTEGRATIONHUB_DATABASE_URL", "DATABASE_URL", "CLOUD_DB_DSN"),
		AutoMigrate: boolEnvDefault(true, "IM_AUTO_MIGRATE"),
		SeedData:    boolEnvDefault(true, "IM_SEED_DATA"),
	}
}

func firstEnv(names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
	}
	return ""
}

func boolEnvDefault(defaultValue bool, name string) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	if value == "" {
		return defaultValue
	}
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}
