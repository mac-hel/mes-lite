package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Host                  string
	Port                  int
	DatabaseURL           string
	MigrationsDir         string
	AuthBootstrapEmail    string
	AuthBootstrapPassword string
	JWTSecret             string
	LogLevel              string
	LogFormat             string
	OTELTracesExporter    string
}

// Load reads configuration from environment variables, applying defaults where
// values are not set.
func Load() Config {
	cfg := Config{
		Host:               "0.0.0.0",
		Port:               9090,
		DatabaseURL:        "postgres://meslite:meslite@localhost:5432/meslite?sslmode=disable",
		MigrationsDir:      "migrations",
		LogLevel:           "info",
		LogFormat:          "json",
		OTELTracesExporter: "none",
	}

	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}

	if p := os.Getenv("PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			cfg.Port = port
		}
	}

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cfg.DatabaseURL = databaseURL
	}

	if migrationsDir := os.Getenv("MIGRATIONS_DIR"); migrationsDir != "" {
		cfg.MigrationsDir = migrationsDir
	}

	if email := os.Getenv("AUTH_BOOTSTRAP_EMAIL"); email != "" {
		cfg.AuthBootstrapEmail = email
	}

	if password := os.Getenv("AUTH_BOOTSTRAP_PASSWORD"); password != "" {
		cfg.AuthBootstrapPassword = password
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWTSecret = secret
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.LogFormat = format
	}

	if exporter := os.Getenv("OTEL_TRACES_EXPORTER"); exporter != "" {
		cfg.OTELTracesExporter = exporter
	}

	return cfg
}

// Addr returns the host:port address string for the HTTP server.
func (c Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
