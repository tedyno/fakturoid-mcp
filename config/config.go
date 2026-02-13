package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Slug         string `json:"slug"`
}

const configDir = "fakturoid-mcp"
const configFile = "config.json"

// Load reads configuration from environment variables with fallback to config file.
// Priority: env > config file
func Load() (*Config, error) {
	cfg := &Config{}

	// Try config file first as base
	if fileCfg, err := loadFromFile(); err == nil {
		cfg = fileCfg
	}

	// Env variables override config file
	if v := os.Getenv("FAKTUROID_CLIENT_ID"); v != "" {
		cfg.ClientID = v
	}
	if v := os.Getenv("FAKTUROID_CLIENT_SECRET"); v != "" {
		cfg.ClientSecret = v
	}
	if v := os.Getenv("FAKTUROID_SLUG"); v != "" {
		cfg.Slug = v
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("FAKTUROID_CLIENT_ID and FAKTUROID_CLIENT_SECRET required (use env variables or ~/.config/%s/%s)", configDir, configFile)
	}
	if cfg.Slug == "" {
		return nil, fmt.Errorf("FAKTUROID_SLUG required (your Fakturoid account slug)")
	}

	return cfg, nil
}

func loadFromFile() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".config", configDir, configFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file %s: %w", path, err)
	}

	return &cfg, nil
}
