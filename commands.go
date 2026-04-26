package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func doInit() error {
	if _, err := os.Stat(configFile); err == nil {
		fmt.Printf("%s already exists. overwrite? [y/N] ", configFile)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			fmt.Println("aborted.")
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(configFile), 0o755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	if err := os.WriteFile(configFile, []byte(initTemplate), 0o644); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	fmt.Printf("created %s\n", configFile)
	fmt.Println("edit it to add your routes, then run: airoute list")
	return nil
}

func doList() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if len(cfg.Routes) == 0 {
		fmt.Println("no routes configured. edit ~/.airoute/config.yaml")
		return nil
	}

	keys := make([]string, 0, len(cfg.Routes))
	for k := range cfg.Routes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Compute column widths.
	colRoute, colURL, colModel, colTier := len("ROUTE"), len("BASE_URL"), len("MODEL"), len("TIER")
	for _, k := range keys {
		r := cfg.Routes[k]
		if len(k) > colRoute {
			colRoute = len(k)
		}
		if len(r.BaseURL) > colURL {
			colURL = len(r.BaseURL)
		}
		if len(r.Model) > colModel {
			colModel = len(r.Model)
		}
		if len(r.Tier) > colTier {
			colTier = len(r.Tier)
		}
	}

	row := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds\n", colRoute, colURL, colModel, colTier)
	fmt.Printf(row, "ROUTE", "BASE_URL", "MODEL", "TIER")
	fmt.Printf(row,
		strings.Repeat("-", colRoute),
		strings.Repeat("-", colURL),
		strings.Repeat("-", colModel),
		strings.Repeat("-", colTier),
	)
	for _, k := range keys {
		r := cfg.Routes[k]
		fmt.Printf(row, k, r.BaseURL, r.Model, r.Tier)
	}
	return nil
}

func doEnv(routeKey string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	route, err := getRoute(cfg, routeKey)
	if err != nil {
		return err
	}
	overrides, err := envOverrides(route)
	if err != nil {
		return err
	}

	// Print in deterministic order.
	keys := []string{
		"ANTHROPIC_BASE_URL",
		"ANTHROPIC_AUTH_TOKEN",
		"ANTHROPIC_API_KEY",
		"ANTHROPIC_MODEL",
		"ANTHROPIC_DEFAULT_OPUS_MODEL",
		"ANTHROPIC_DEFAULT_SONNET_MODEL",
		"ANTHROPIC_DEFAULT_HAIKU_MODEL",
	}
	for _, k := range keys {
		if val, ok := overrides[k]; ok {
			// Single-quote the value and escape any embedded single quotes.
			escaped := strings.ReplaceAll(val, "'", "'\\''")
			fmt.Printf("export %s='%s'\n", k, escaped)
		}
	}
	return nil
}

func doRun(routeKey string, passthrough []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	route, err := getRoute(cfg, routeKey)
	if err != nil {
		return err
	}
	env, err := buildEnv(route, os.Environ())
	if err != nil {
		return err
	}

	binary := "claude"
	args := []string{"claude"}
	if len(passthrough) > 0 {
		binary = passthrough[0]
		args = passthrough
	}

	path, err := exec.LookPath(binary)
	if err != nil {
		if binary == "claude" {
			return fmt.Errorf("'claude' not found in PATH\ninstall: npm install -g @anthropic-ai/claude-code")
		}
		return fmt.Errorf("%q not found in PATH", binary)
	}

	// Replace this process — no parent wrapper left behind.
	return syscall.Exec(path, args, env)
}
