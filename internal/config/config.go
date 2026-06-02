package config

import (
	"os"
	"strings"
)

// Config holds all application configuration.
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// CORS
	CORSOrigins []string

	// Team Engine
	TeamEngineEnabled  bool
	TeamEngineInterval string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://tnotes:tnotes@localhost:5432/tnotes_teams?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "tnotes-dev-secret-change-in-production"),
		JWTExpiration:      getEnv("JWT_EXPIRATION", "24h"),
		CORSOrigins:        strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000"), ","),
		TeamEngineEnabled:  getEnv("TEAM_ENGINE_ENABLED", "true") == "true",
		TeamEngineInterval: getEnv("TEAM_ENGINE_INTERVAL", "30s"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
