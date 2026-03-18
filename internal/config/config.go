package config

import (
	"fmt"
	"os"
)

type Config struct {
	GRPCPort     string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresSSL  string
}

func Load() Config {
	cfg := Config{
		GRPCPort:     getEnv("GRPC_PORT", "50051"),
		PostgresHost: getEnv("POSTGRES_HOST", "127.0.0.1"),
		PostgresPort: getEnv("POSTGRES_PORT", "5434"),
		PostgresUser: getEnv("POSTGRES_USER", "postgres"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:   getEnv("POSTGRES_DB", "shipments"),
		PostgresSSL:  getEnv("POSTGRES_SSLMODE", "disable"),
	}

	return cfg
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser,
		c.PostgresPass,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
		c.PostgresSSL,
	)
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}