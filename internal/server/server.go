package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"

	"github.com/mac-hel/mes-lite/internal/auth"
	"github.com/mac-hel/mes-lite/internal/csvimport"
	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/machines"
	"github.com/mac-hel/mes-lite/internal/orders"
	"github.com/mac-hel/mes-lite/internal/platform/config"
	"github.com/mac-hel/mes-lite/internal/platform/jobs"
	"github.com/mac-hel/mes-lite/internal/platform/logging"
	"github.com/mac-hel/mes-lite/internal/platform/metrics"
	"github.com/mac-hel/mes-lite/internal/platform/tracing"
	"github.com/mac-hel/mes-lite/internal/platform/version"
	"github.com/mac-hel/mes-lite/internal/production"
	"github.com/mac-hel/mes-lite/internal/products"
	"github.com/mac-hel/mes-lite/internal/reporting"
)

// Server wraps the Fuego HTTP server with application-specific configuration.
type Server struct {
	*fuego.Server
	cfg            config.Config
	logger         *slog.Logger
	readinessCheck func(context.Context) error
}

// New creates a new Server with the given configuration and registers routes.
func New(cfg config.Config, authHandler *auth.Handler, authMiddleware *auth.Middleware, empHandler *employees.Handler, prodHandler *products.Handler, productionHandler *production.Handler, ordersHandler *orders.Handler, reportingHandler *reporting.Handler, csvImportHandlers ...*csvimport.Handler) *Server {
	return NewWithLogger(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), authHandler, authMiddleware, empHandler, prodHandler, productionHandler, ordersHandler, reportingHandler, csvImportHandlers...)
}

// NewWithLogger creates a Server that logs lifecycle and request records through logger.
func NewWithLogger(cfg config.Config, logger *slog.Logger, authHandler *auth.Handler, authMiddleware *auth.Middleware, empHandler *employees.Handler, prodHandler *products.Handler, productionHandler *production.Handler, ordersHandler *orders.Handler, reportingHandler *reporting.Handler, csvImportHandlers ...*csvimport.Handler) *Server {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	httpMetrics := metrics.NewHTTPMetrics()

	fuegoServer := fuego.NewServer(
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
	s := &Server{Server: fuegoServer, cfg: cfg, logger: logger}
	fuego.GetStd(s.Server, "/metrics", httpMetrics.Handler().ServeHTTP, fuego.OptionHide())
	fuego.Use(s.Server, httpMetrics.Middleware)
	fuego.Use(s.Server, tracing.Middleware(nil))
	fuego.Use(s.Server, logging.RequestLogger(logger))
	requireBearer := fuego.OptionSecurity(openapi3.NewSecurityRequirement().Authenticate("bearerAuth"))
	adminOnly := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin)))
	readMasterData := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader)))
	manageProducts := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager)))
	manageOrders := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager)))
	progressOrders := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader)))
	registerProduction := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader, auth.RoleWorker)))
	readReports := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader)))
	importProduction := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager)))
	reviewProduction := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager, auth.RoleLeader)))

	fuego.Get(s.Server, "/ready", s.readyHandler)
	fuego.Get(s.Server, "/health", healthHandler)
	fuego.Get(s.Server, "/version", versionHandler)
	fuego.Post(s.Server, "/auth/login", authHandler.Login)

	fuego.Post(s.Server, "/employees", empHandler.Create, adminOnly)
	fuego.Get(s.Server, "/employees", empHandler.List, readMasterData)
	fuego.Put(s.Server, "/employees/{id}", empHandler.Update, adminOnly)
	fuego.Put(s.Server, "/employees/{id}/deactivate", empHandler.Deactivate, adminOnly)

	fuego.Post(s.Server, "/products", prodHandler.Create, manageProducts)
	fuego.Get(s.Server, "/products", prodHandler.List, readMasterData)
	fuego.Get(s.Server, "/products/search", prodHandler.Search, readMasterData)
	fuego.Put(s.Server, "/products/{sku}", prodHandler.Update, manageProducts)
	fuego.Put(s.Server, "/products/{sku}/deactivate", prodHandler.Deactivate, manageProducts)

	fuego.Post(s.Server, "/production-entries", productionHandler.Register, registerProduction)
	fuego.Get(s.Server, "/production-entries", productionHandler.List, reviewProduction)
	fuego.Post(s.Server, "/production-entries/{id}/corrections", productionHandler.Correct, reviewProduction)
	fuego.Get(s.Server, "/production-entries/{id}/corrections", productionHandler.ListCorrections, reviewProduction)

	fuego.Post(s.Server, "/production-orders", ordersHandler.Create, manageOrders)
	fuego.Get(s.Server, "/production-orders/{id}", ordersHandler.Get, readMasterData)
	fuego.Post(s.Server, "/production-orders/{id}/assignments", ordersHandler.AssignEmployee, manageOrders)
	fuego.Put(s.Server, "/production-orders/{id}/release", ordersHandler.Release, progressOrders)
	fuego.Put(s.Server, "/production-orders/{id}/start", ordersHandler.Start, progressOrders)
	fuego.Put(s.Server, "/production-orders/{id}/complete", ordersHandler.Complete, progressOrders)
	fuego.Put(s.Server, "/production-orders/{id}/cancel", ordersHandler.Cancel, progressOrders)

	fuego.Get(s.Server, "/reports/daily-production", reportingHandler.DailyProduction, readReports)
	fuego.Get(s.Server, "/reports/daily-employee-production", reportingHandler.DailyEmployeeProduction, readReports)
	fuego.Get(s.Server, "/reports/employee-productivity", reportingHandler.EmployeeProductivity, readReports)
	fuego.Get(s.Server, "/reports/employee-productivity/products", reportingHandler.EmployeeProductivityProducts, readReports)
	fuego.Get(s.Server, "/reports/product-statistics", reportingHandler.ProductStatistics, readReports)

	if len(csvImportHandlers) > 0 && csvImportHandlers[0] != nil {
		fuego.Post(s.Server, "/imports/production-entries", csvImportHandlers[0].ImportProductionEntries, importProduction)
		if csvImportHandlers[0].AsyncEnabled() {
			fuego.Post(s.Server, "/imports/production-entries/jobs", csvImportHandlers[0].ImportProductionEntriesAsync, importProduction)
		}
	}

	return s
}

// SetReadinessCheck installs the dependency check used by /ready.
func (s *Server) SetReadinessCheck(check func(context.Context) error) {
	s.readinessCheck = check
}

// RegisterJobRoutes registers operational background-job routes.
func RegisterJobRoutes(s *Server, authMiddleware *auth.Middleware, jobsHandler *jobs.HTTPHandler) {
	if jobsHandler == nil {
		return
	}

	requireBearer := fuego.OptionSecurity(openapi3.NewSecurityRequirement().Authenticate("bearerAuth"))
	manageJobs := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager)))

	fuego.Get(s.Server, "/jobs/{id}", jobsHandler.Get, manageJobs)
	fuego.Put(s.Server, "/jobs/{id}/cancel", jobsHandler.Cancel, manageJobs)
}

// RegisterMachineRoutes registers fake machine integration routes.
func RegisterMachineRoutes(s *Server, authMiddleware *auth.Middleware, machineHandler *machines.Handler) {
	if machineHandler == nil {
		return
	}

	requireBearer := fuego.OptionSecurity(openapi3.NewSecurityRequirement().Authenticate("bearerAuth"))
	manageMachines := fuego.GroupOptions(requireBearer, fuego.OptionMiddleware(authMiddleware.Authenticate, authMiddleware.RequireRole(auth.RoleAdmin, auth.RoleManager)))

	fuego.Get(s.Server, "/machines/events/stats", machineHandler.Stats, manageMachines)
	fuego.Post(s.Server, "/machines/{machineId}/events", machineHandler.CreateEvent, manageMachines)
}

type readyResponse struct {
	Status string `json:"status"`
}

func (s *Server) readyHandler(c fuego.ContextNoBody) (readyResponse, error) {
	if s.readinessCheck == nil {
		return readyResponse{Status: "ready"}, nil
	}
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()
	if err := s.readinessCheck(ctx); err != nil {
		s.logger.WarnContext(c.Context(), "readiness check failed", "err", err)
		return readyResponse{}, fuego.HTTPError{Err: err, Status: http.StatusServiceUnavailable, Detail: "service is not ready"}
	}
	return readyResponse{Status: "ready"}, nil
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
		s.logger.Info("server starting", "addr", s.cfg.Addr())
		if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server error", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	s.logger.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Shutdown(shutdownCtx)
}
