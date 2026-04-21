package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        int
	BaseURL     string
	DatabaseURL string
}

func Load() (Config, error) {
	cfg := Config{
		Port:        envInt("PORT", 8080),
		BaseURL:     strings.TrimRight(envString("BASE_URL", "http://localhost:8080"), "/"),
		DatabaseURL: envString("DATABASE_URL", ""),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func envString(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
