package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8082"),
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://korp_billing:korp_billing_password@localhost:5433/billing_db?sslmode=disable",
		),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func (c Config) Address() string {
	return fmt.Sprintf(":%s", c.Port)
}
