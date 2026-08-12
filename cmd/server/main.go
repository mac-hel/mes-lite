package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/config"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	os.Exit(run())
}

func run() int {
	config.LoadDotEnv(".env")
	cfg := config.Load()
	ctx := context.Background()

	poolCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(poolCtx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("open database pool", "err", err)
		return 1
	}
	defer db.Close()

	if err := db.Ping(poolCtx); err != nil {
		slog.Error("ping database", "err", err)
		return 1
	}
	if cfg.JWTSecret == "" {
		slog.Error("JWT_SECRET is required")
		return 1
	}
	tokens, err := auth.NewTokenManager(cfg.JWTSecret)
	if err != nil {
		slog.Error("create token manager", "err", err)
		return 1
	}

	authStore := auth.NewInMemoryStore()
	if cfg.AuthBootstrapEmail != "" || cfg.AuthBootstrapPassword != "" {
		if cfg.AuthBootstrapEmail == "" || cfg.AuthBootstrapPassword == "" {
			slog.Error("AUTH_BOOTSTRAP_EMAIL and AUTH_BOOTSTRAP_PASSWORD must be set together")
			return 1
		}
		user, err := auth.NewUser("bootstrap-admin", cfg.AuthBootstrapEmail, cfg.AuthBootstrapPassword, auth.RoleAdmin)
		if err != nil {
			slog.Error("create bootstrap auth user", "err", err)
			return 1
		}
		if err := authStore.Save(ctx, user); err != nil {
			slog.Error("save bootstrap auth user", "err", err)
			return 1
		}
	}
	authHandler := auth.NewHandler(auth.NewService(authStore, tokens))

	empStore := employees.NewPostgresStore(db)
	empHandler := employees.NewHandler(empStore)

	prodStore := products.NewPostgresStore(db)
	prodHandler := products.NewHandler(prodStore)

	productionStore := production.NewPostgresStore(db)
	productionService := production.NewService(productionStore, empStore, prodStore)
	productionHandler := production.NewHandler(productionService)

	srv := server.New(cfg, authHandler, empHandler, prodHandler, productionHandler)

	if err := srv.Start(ctx); err != nil {
		slog.Error("server shutdown with error", "err", err)
		return 1
	}

	return 0
}
