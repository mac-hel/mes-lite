package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const jwtSecretMinLength = 32

var (
	// ErrInvalidConfig means required runtime configuration is missing or invalid.
	ErrInvalidConfig = errors.New("invalid config")
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

// ValidateServer checks configuration required to run the HTTP API process.
func (c Config) ValidateServer() error {
	if err := c.validateCommon(); err != nil {
		return err
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required: %w", ErrInvalidConfig)
	}
	if len(c.JWTSecret) < jwtSecretMinLength {
		return fmt.Errorf("JWT_SECRET must be at least %d characters: %w", jwtSecretMinLength, ErrInvalidConfig)
	}
	if (c.AuthBootstrapEmail == "") != (c.AuthBootstrapPassword == "") {
		return fmt.Errorf("AUTH_BOOTSTRAP_EMAIL and AUTH_BOOTSTRAP_PASSWORD must be set together: %w", ErrInvalidConfig)
	}
	return nil
}

// ValidateMigrate checks configuration required to run database migrations.
func (c Config) ValidateMigrate() error {
	return c.validateCommon()
}

func (c Config) validateCommon() error {
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("HOST is required: %w", ErrInvalidConfig)
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535: %w", ErrInvalidConfig)
	}
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("DATABASE_URL is required: %w", ErrInvalidConfig)
	}
	if strings.TrimSpace(c.MigrationsDir) == "" {
		return fmt.Errorf("MIGRATIONS_DIR is required: %w", ErrInvalidConfig)
	}
	return nil
}

// Addr returns the host:port address string for the HTTP server.
func (c Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
