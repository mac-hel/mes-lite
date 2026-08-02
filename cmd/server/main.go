package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mac-hel/mes-lite/internal/config"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	config.LoadDotEnv(".env")
	cfg := config.Load()

	empStore := employees.NewInMemoryStore()
	empHandler := employees.NewHandler(empStore)

	prodStore := products.NewInMemoryStore()
	prodHandler := products.NewHandler(prodStore)

	srv := server.New(cfg, empHandler, prodHandler)

	if err := srv.Start(context.Background()); err != nil {
		slog.Error("server shutdown with error", "err", err)
		os.Exit(1)
	}
}
