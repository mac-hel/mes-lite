package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	_ = os.Unsetenv("HOST")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Unsetenv("MIGRATIONS_DIR")
	_ = os.Unsetenv("AUTH_BOOTSTRAP_EMAIL")
	_ = os.Unsetenv("AUTH_BOOTSTRAP_PASSWORD")
	_ = os.Unsetenv("JWT_SECRET")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("LOG_FORMAT")
	_ = os.Unsetenv("OTEL_TRACES_EXPORTER")

	cfg := Load()

	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected default host 0.0.0.0, got %s", cfg.Host)
	}

	if cfg.Port != 9090 {
		t.Errorf("expected default port 9090, got %d", cfg.Port)
	}

	if cfg.DatabaseURL != "postgres://meslite:meslite@localhost:5432/meslite?sslmode=disable" {
		t.Errorf("expected default database URL, got %s", cfg.DatabaseURL)
	}

	if cfg.MigrationsDir != "migrations" {
		t.Errorf("expected default migrations dir migrations, got %s", cfg.MigrationsDir)
	}

	if cfg.AuthBootstrapEmail != "" {
		t.Errorf("expected empty auth bootstrap email, got %s", cfg.AuthBootstrapEmail)
	}

	if cfg.AuthBootstrapPassword != "" {
		t.Error("expected empty auth bootstrap password")
	}

	if cfg.JWTSecret != "" {
		t.Error("expected empty jwt secret")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("expected default log level info, got %s", cfg.LogLevel)
	}

	if cfg.LogFormat != "json" {
		t.Errorf("expected default log format json, got %s", cfg.LogFormat)
	}

	if cfg.OTELTracesExporter != "none" {
		t.Errorf("expected default trace exporter none, got %s", cfg.OTELTracesExporter)
	}
}

func TestLoadFromEnv(t *testing.T) {
	_ = os.Setenv("HOST", "127.0.0.1")
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("DATABASE_URL", "postgres://example")
	_ = os.Setenv("MIGRATIONS_DIR", "db/migrations")
	_ = os.Setenv("AUTH_BOOTSTRAP_EMAIL", "admin@example.com")
	_ = os.Setenv("AUTH_BOOTSTRAP_PASSWORD", "secret")
	_ = os.Setenv("JWT_SECRET", "jwt-secret-with-at-least-32-characters")
	_ = os.Setenv("LOG_LEVEL", "debug")
	_ = os.Setenv("LOG_FORMAT", "text")
	_ = os.Setenv("OTEL_TRACES_EXPORTER", "stdout")
	defer func() {
		_ = os.Unsetenv("HOST")
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("MIGRATIONS_DIR")
		_ = os.Unsetenv("AUTH_BOOTSTRAP_EMAIL")
		_ = os.Unsetenv("AUTH_BOOTSTRAP_PASSWORD")
		_ = os.Unsetenv("JWT_SECRET")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("LOG_FORMAT")
		_ = os.Unsetenv("OTEL_TRACES_EXPORTER")
	}()

	cfg := Load()

	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Host)
	}

	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Errorf("expected database URL from env, got %s", cfg.DatabaseURL)
	}

	if cfg.MigrationsDir != "db/migrations" {
		t.Errorf("expected migrations dir from env, got %s", cfg.MigrationsDir)
	}

	if cfg.AuthBootstrapEmail != "admin@example.com" {
		t.Errorf("expected auth bootstrap email from env, got %s", cfg.AuthBootstrapEmail)
	}

	if cfg.AuthBootstrapPassword != "secret" {
		t.Error("expected auth bootstrap password from env")
	}

	if cfg.JWTSecret != "jwt-secret-with-at-least-32-characters" {
		t.Error("expected jwt secret from env")
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("expected log level from env, got %s", cfg.LogLevel)
	}

	if cfg.LogFormat != "text" {
		t.Errorf("expected log format from env, got %s", cfg.LogFormat)
	}

	if cfg.OTELTracesExporter != "stdout" {
		t.Errorf("expected trace exporter from env, got %s", cfg.OTELTracesExporter)
	}
}

func TestAddr(t *testing.T) {
	cfg := Config{Host: "0.0.0.0", Port: 9090}

	if got := cfg.Addr(); got != "0.0.0.0:9090" {
		t.Errorf("expected 0.0.0.0:9090, got %s", got)
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("TEST_DOTENV=from_file"), 0644); err != nil {
		t.Fatal(err)
	}

	_ = os.Unsetenv("TEST_DOTENV")
	defer func() { _ = os.Unsetenv("TEST_DOTENV") }()

	LoadDotEnv(path)

	if got := os.Getenv("TEST_DOTENV"); got != "from_file" {
		t.Errorf("expected from_file, got %s", got)
	}
}

func TestLoadDotEnvMissingFileIsNoop(t *testing.T) {
	LoadDotEnv("/tmp/nonexistent_dotenv_file_12345")
}
