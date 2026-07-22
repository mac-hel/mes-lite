package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Host string
	Port int
}

// Load reads configuration from environment variables, applying defaults where
// values are not set.
func Load() Config {
	cfg := Config{
		Host: "0.0.0.0",
		Port: 8080,
	}

	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}

	if p := os.Getenv("PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			cfg.Port = port
		}
	}

	return cfg
}

// Addr returns the host:port address string for the HTTP server.
func (c Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
