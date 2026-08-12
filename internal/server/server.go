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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/config"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/version"
)

// Server wraps the Fuego HTTP server with application-specific configuration.
type Server struct {
	*fuego.Server
	cfg config.Config
}

// New creates a new Server with the given configuration and registers routes.
func New(cfg config.Config, authHandler *auth.Handler, authMiddleware *auth.Middleware, empHandler *employees.Handler, prodHandler *products.Handler, productionHandler *production.Handler) *Server {
	s := fuego.NewServer(
		fuego.WithAddr(cfg.Addr()),
		fuego.WithSecurity(openapi3.SecuritySchemes{
			"bearerAuth": &openapi3.SecuritySchemeRef{
				Value: openapi3.NewSecurityScheme().
					WithType("http").
					WithScheme("bearer").
					WithBearerFormat("JWT").
					WithDescription("JWT access token from POST /auth/login"),
			},
		}),
	)
	requireBearer := fuego.OptionSecurity(openapi3.NewSecurityRequirement().Authenticate("bearerAuth"))
	adminOnly := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin)))
	readMasterData := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader)))
	manageProducts := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager)))
	registerProduction := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader, auth.RoleWorker)))

	fuego.Get(s, "/ready", readyHandler)
	fuego.Get(s, "/health", healthHandler)
	fuego.Get(s, "/version", versionHandler)
	fuego.Post(s, "/auth/login", authHandler.Login)

	fuego.Post(s, "/employees", empHandler.Create, adminOnly)
	fuego.Get(s, "/employees", empHandler.List, readMasterData)
	fuego.Put(s, "/employees/{id}", empHandler.Update, adminOnly)
	fuego.Put(s, "/employees/{id}/deactivate", empHandler.Deactivate, adminOnly)

	fuego.Post(s, "/products", prodHandler.Create, manageProducts)
	fuego.Get(s, "/products", prodHandler.List, readMasterData)
	fuego.Get(s, "/products/search", prodHandler.Search, readMasterData)
	fuego.Put(s, "/products/{sku}", prodHandler.Update, manageProducts)
	fuego.Put(s, "/products/{sku}/deactivate", prodHandler.Deactivate, manageProducts)

	fuego.Post(s, "/production-entries", productionHandler.Register, registerProduction)

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
