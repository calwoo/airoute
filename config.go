package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var configFile = filepath.Join(os.Getenv("HOME"), ".airoute", "config.yaml")

var validTiers = map[string]bool{"opus": true, "sonnet": true, "haiku": true}

type Route struct {
	BaseURL     string `yaml:"base_url"`
	Model       string `yaml:"model"`
	APIKeyEnv   string `yaml:"api_key_env"`
	ClearAPIKey bool   `yaml:"clear_api_key"`
	Tier        string `yaml:"tier"`
}

type Config struct {
	Routes map[string]Route `yaml:"routes"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found — run: airoute init")
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}
	if cfg.Routes == nil {
		return nil, fmt.Errorf("config must have a 'routes' mapping")
	}
	return &cfg, nil
}

func getRoute(cfg *Config, key string) (Route, error) {
	route, ok := cfg.Routes[key]
	if !ok {
		keys := make([]string, 0, len(cfg.Routes))
		for k := range cfg.Routes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return Route{}, fmt.Errorf("route %q not found. available: %s", key, strings.Join(keys, ", "))
	}
	if route.BaseURL == "" {
		return Route{}, fmt.Errorf("route %q: missing required field: base_url", key)
	}
	if route.Model == "" {
		return Route{}, fmt.Errorf("route %q: missing required field: model", key)
	}
	if route.Tier != "" && !validTiers[route.Tier] {
		return Route{}, fmt.Errorf("route %q: tier must be opus, sonnet, or haiku", key)
	}
	return route, nil
}
