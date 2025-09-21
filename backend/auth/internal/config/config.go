package config

import (
	"fmt"
	"net/url"
	"os"
)

// Config represents runtime configuration for the auth service.
type Config struct {
	HTTPAddr string
	Issuer   string
}

// Load builds configuration from environment variables with sane defaults.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr: getEnv("AUTH_HTTP_ADDR", ":8080"),
		Issuer:   getEnv("AUTH_TOKEN_ISSUER", "https://auth.tchat.dev"),
	}

	if _, err := url.ParseRequestURI(cfg.Issuer); err != nil {
		return Config{}, fmt.Errorf("invalid issuer %q: %w", cfg.Issuer, err)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
