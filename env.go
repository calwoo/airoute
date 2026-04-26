package main

import (
	"fmt"
	"os"
	"strings"
)

var tierEnvVar = map[string]string{
	"opus":   "ANTHROPIC_DEFAULT_OPUS_MODEL",
	"sonnet": "ANTHROPIC_DEFAULT_SONNET_MODEL",
	"haiku":  "ANTHROPIC_DEFAULT_HAIKU_MODEL",
}

// buildEnv returns a copy of environ with the route's vars injected.
// environ should be os.Environ().
func buildEnv(route Route, environ []string) ([]string, error) {
	overrides := map[string]string{
		"ANTHROPIC_BASE_URL": route.BaseURL,
		"ANTHROPIC_MODEL":    route.Model,
	}

	if route.APIKeyEnv != "" {
		val := os.Getenv(route.APIKeyEnv)
		if val == "" {
			return nil, fmt.Errorf(
				"env var %q is not set in your shell\nrun: export %s=<your-key>",
				route.APIKeyEnv, route.APIKeyEnv,
			)
		}
		overrides["ANTHROPIC_AUTH_TOKEN"] = val
	}

	if route.ClearAPIKey {
		overrides["ANTHROPIC_API_KEY"] = ""
	}

	if route.Tier != "" {
		overrides[tierEnvVar[route.Tier]] = route.Model
	}

	// Rebuild environ, replacing any keys we override.
	seen := make(map[string]bool, len(overrides))
	result := make([]string, 0, len(environ)+len(overrides))
	for _, entry := range environ {
		key, _, _ := strings.Cut(entry, "=")
		if val, ok := overrides[key]; ok {
			result = append(result, key+"="+val)
			seen[key] = true
		} else {
			result = append(result, entry)
		}
	}
	// Append any override keys that weren't already in environ.
	for key, val := range overrides {
		if !seen[key] {
			result = append(result, key+"="+val)
		}
	}
	return result, nil
}

// envOverrides returns only the keys that airoute manages, for use by doEnv.
func envOverrides(route Route) (map[string]string, error) {
	env, err := buildEnv(route, nil)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(env))
	for _, entry := range env {
		key, val, _ := strings.Cut(entry, "=")
		m[key] = val
	}
	return m, nil
}
