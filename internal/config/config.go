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
}

// Load reads configuration from environment variables, applying defaults where
// values are not set.
func Load() Config {
	cfg := Config{
		Host:          "0.0.0.0",
		Port:          9090,
		DatabaseURL:   "postgres://meslite:meslite@localhost:5432/meslite?sslmode=disable",
		MigrationsDir: "migrations",
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

	return cfg
}

// Addr returns the host:port address string for the HTTP server.
func (c Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
