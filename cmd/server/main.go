// Package main is the entry point for the MES-Lite server application.
// It initializes the server, sets up the database connection, configures
// authentication, and starts the HTTP server with the appropriate handlers for
// various endpoints.
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/csvimport"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/orders"
	"github.com/mac-hel/mes-lite/internal/platform/config"
	"github.com/mac-hel/mes-lite/internal/platform/jobs"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/reporting"
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

	authStore := auth.NewPostgresStore(db)
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
		if err := authStore.EnsureBootstrapAdmin(ctx, user); err != nil {
			slog.Error("save bootstrap auth user", "err", err)
			return 1
		}
	}
	authHandler := auth.NewHandler(auth.NewService(authStore, tokens))
	authMiddleware := auth.NewMiddleware(tokens)

	empStore := employees.NewPostgresStore(db)
	empHandler := employees.NewHandler(empStore)

	prodStore := products.NewPostgresStore(db)
	prodHandler := products.NewHandler(prodStore)

	productionStore := production.NewPostgresStore(db)
	productionService := production.NewService(productionStore, productionStore, empStore, prodStore)
	productionHandler := production.NewHandler(productionService, productionService)

	ordersStore := orders.NewPostgresStore(db)
	ordersService := orders.NewService(ordersStore, empStore, prodStore)
	ordersHandler := orders.NewHandler(ordersService)

	reportingStore := reporting.NewPostgresStore(db)
	reportingHandler := reporting.NewHandler(reportingStore)

	csvImportStore := csvimport.NewPostgresStore(db)
	csvImportService := csvimport.NewService(csvImportStore)
	jobQueue := jobs.NewQueue(jobs.DefaultQueueCapacity)
	jobWorkers, err := jobs.NewWorkerPool(jobQueue, jobs.DefaultWorkerCount, map[jobs.Type]jobs.Handler{
		jobs.TypeProductionEntryImport: csvimport.NewProductionEntriesJobHandler(csvImportService, jobQueue),
	})
	if err != nil {
		slog.Error("create background worker pool", "err", err)
		return 1
	}
	csvImportAsyncService := csvimport.NewAsyncService(jobQueue, "")
	csvImportHandler := csvimport.NewHandlerWithAsync(csvImportService, csvImportAsyncService)
	jobWorkers.Start(ctx)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := jobWorkers.Stop(shutdownCtx); err != nil {
			slog.Error("stop background worker pool", "err", err)
		}
	}()
	jobsHandler := jobs.NewHTTPHandler(jobQueue, jobWorkers)

	srv := server.New(cfg, authHandler, authMiddleware, empHandler, prodHandler, productionHandler, ordersHandler, reportingHandler, csvImportHandler)
	server.RegisterJobRoutes(srv, authMiddleware, jobsHandler)

	if err := srv.Start(ctx); err != nil {
		slog.Error("server shutdown with error", "err", err)
		return 1
	}

	return 0
}
