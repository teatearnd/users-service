package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	DatabaseUrl string
	JWTSecret   string
	Port        string
}

func LoadConfig() *Config {
	cfg := &Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseUrl == "" {
		log.Fatalf("Config error: database is required")
	}

	if !strings.HasPrefix(cfg.DatabaseUrl, "postgres://") && !strings.HasPrefix(cfg.DatabaseUrl, "host=") {
		log.Fatalf("Config error: DATABASE_URL must be a valid connection string")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("Config error: JWT_SECRET is required")
	}

	if cfg.Port == "" {
		log.Fatalf("Config error: PORT not provided")
	}
	return cfg
}
