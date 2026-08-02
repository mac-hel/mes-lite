package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-fuego/fuego"

	"github.com/mac-hel/mes-lite/internal/config"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/version"
)

// Server wraps the Fuego HTTP server with application-specific configuration.
type Server struct {
	*fuego.Server
	cfg config.Config
}

// New creates a new Server with the given configuration and registers routes.
func New(cfg config.Config, empHandler *employees.Handler, prodHandler *products.Handler) *Server {
	s := fuego.NewServer(
		fuego.WithAddr(cfg.Addr()),
	)

	fuego.Get(s, "/ready", readyHandler)
	fuego.Get(s, "/health", healthHandler)
	fuego.Get(s, "/version", versionHandler)

	fuego.Post(s, "/employees", empHandler.Create)
	fuego.Get(s, "/employees", empHandler.List)
	fuego.Put(s, "/employees/{id}", empHandler.Update)
	fuego.Put(s, "/employees/{id}/deactivate", empHandler.Deactivate)

	fuego.Post(s, "/products", prodHandler.Create)
	fuego.Get(s, "/products", prodHandler.List)
	fuego.Get(s, "/products/search", prodHandler.Search)
	fuego.Put(s, "/products/{sku}", prodHandler.Update)
	fuego.Put(s, "/products/{sku}/deactivate", prodHandler.Deactivate)

	return &Server{Server: s, cfg: cfg}
}

type readyResponse struct{}

func readyHandler(c fuego.ContextNoBody) (readyResponse, error) {
	return readyResponse{}, nil
}

type healthResponse struct {
	Status string `json:"status"`
}

func healthHandler(c fuego.ContextNoBody) (healthResponse, error) {
	return healthResponse{Status: "ok"}, nil
}

type versionResponse struct {
	Version string `json:"version"`
}

func versionHandler(c fuego.ContextNoBody) (versionResponse, error) {
	return versionResponse{Version: version.String()}, nil
}

// Start runs the server and blocks until SIGINT or SIGTERM is received,
// then performs a graceful shutdown with a 10-second timeout.
func (s *Server) Start(ctx context.Context) error {
	// os provides abstraction for platform-specific signals:
	// 	os.Kill      = syscall.SIGKILL, force exit (`kill -9`)
	// 	os.Interrupt = syscall.SIGINT, interrupt (Ctrl+C)
	// 	syscall.SIGTERM terminate (`kill`), not all systems support it
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server starting", "addr", s.cfg.Addr())
		if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Shutdown(shutdownCtx)
}
