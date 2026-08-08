// package docs
// `go doc ./cmd/server` displays it
package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/mac-hel/mes-lite/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	os.Exit(run())
}

func run() int {
	config.LoadDotEnv(".env")
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		slog.Error("open database", "err", err)
		return 1
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close database", "err", err)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("ping database", "err", err)
		return 1
	}

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("set migration dialect", "err", err)
		return 1
	}

	slog.Info("running migrations", "dir", cfg.MigrationsDir)
	if err := goose.Up(db, cfg.MigrationsDir); err != nil {
		slog.Error("run migrations", "err", err)
		return 1
	}
	slog.Info("migrations complete")
	return 0
}
