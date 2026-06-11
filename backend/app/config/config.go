package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Simaky/go-github-tracker/backend/consts"
)

const defaultListen = ":12010"

// Config is the typed view of everything the service reads from its environment.
type Config struct {
	Storage Storage
	GitHub  GitHub
	Listen  string `env:"GOGITHUBTRACKER_LISTEN"`
}

// Storage holds database connection settings.
type Storage struct {
	DSN string `env:"GOGITHUBTRACKER_DB_DSN"`
}

// GitHub holds settings for the GitHub API client.
type GitHub struct {
	// Token is optional; when set it raises the API rate limit.
	Token string `env:"GITHUB_TOKEN"`
}

// LoadConfig resolves configuration from the environment:
//  1. If a .env file is present (local dev), load it — values it sets do NOT
//     override variables already exported in the real environment.
//  2. Decode env vars into the typed Config via struct tags.
//  3. Apply defaults, then validate (fatal on missing required fields).
func LoadConfig() Config {
	// Load .env if present. Missing file is fine (e.g. in Docker, where compose
	// injects the variables directly). Real env vars always win.
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parsing env config: %s", err)
	}

	applyDefaults(&cfg)
	validateConfig(cfg)
	return cfg
}

func applyDefaults(cfg *Config) {
	if cfg.Listen == "" {
		cfg.Listen = defaultListen
	}
}

func validateConfig(cfg Config) {
	if cfg.Storage.DSN == "" {
		log.Fatalf("config: %s_DB_DSN is required", strings.ToUpper(consts.ServiceName))
	}
}
