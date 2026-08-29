package config

import "github.com/joho/godotenv"

// LoadDotEnv loads environment variables from a file (typically .env).
func LoadDotEnv(path string) {
	_ = godotenv.Load(path)
}
