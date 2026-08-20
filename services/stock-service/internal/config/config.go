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
		Port: getEnv("PORT", "8081"),
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://korp_stock:korp_stock_password@localhost:5432/stock_db?sslmode=disable",
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