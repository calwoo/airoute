package main

const initTemplate = `# ~/.airoute/config.yaml
# airoute configuration — routes map friendly names to API endpoints + models.
#
# Usage:
#   airoute list                          show all configured routes
#   airoute <route-key>                   start claude with that route's env vars
#   airoute <route-key> -- <command>      run any command with the env vars set
#   eval $(airoute env <route-key>)       export vars into your current shell
#
# NOTE: "airoute <route> | claude" does NOT work — pipes cannot inject env vars
# into a child process. Use "airoute <route>" directly instead.
#
# Fields per route:
#   base_url     (required) API endpoint URL
#   model        (required) model identifier string for the provider
#   api_key_env  name of the shell env var holding your API key (or "" for none)
#   clear_api_key  set ANTHROPIC_API_KEY="" — required for most third-party providers
#   tier         which Claude tier this maps to: opus | sonnet | haiku
#                sets ANTHROPIC_DEFAULT_{TIER}_MODEL so subagents use the right model

routes:

  # --- Direct Anthropic ---
  # Requires: export ANTHROPIC_API_KEY=sk-ant-...
  anthropic:
    base_url: https://api.anthropic.com
    model: claude-sonnet-4-6
    api_key_env: ANTHROPIC_API_KEY
    tier: sonnet

  # --- OpenRouter ---
  # Requires: export OPENROUTER_API_KEY=sk-or-...
  # clear_api_key is required — OpenRouter expects ANTHROPIC_API_KEY to be empty.
  openrouter-sonnet:
    base_url: https://openrouter.ai/api
    model: anthropic/claude-sonnet-4-6
    api_key_env: OPENROUTER_API_KEY
    clear_api_key: true
    tier: sonnet

  openrouter-opus:
    base_url: https://openrouter.ai/api
    model: anthropic/claude-opus-4-5
    api_key_env: OPENROUTER_API_KEY
    clear_api_key: true
    tier: opus

  # --- DeepInfra ---
  # Requires: export DEEPINFRA_TOKEN=...
  # Supports any model available at https://deepinfra.com/models
  deepinfra-qwen:
    base_url: https://api.deepinfra.com/anthropic
    model: Qwen/Qwen3-30B-A3B
    api_key_env: DEEPINFRA_TOKEN
    clear_api_key: true
    tier: sonnet

  # --- Fireworks ---
  # Requires: export FIREWORKS_API_KEY=...
  fireworks-llama:
    base_url: https://api.fireworks.ai/inference
    model: accounts/fireworks/models/llama-v3p1-70b-instruct
    api_key_env: FIREWORKS_API_KEY
    clear_api_key: true
    tier: sonnet

  # --- Local (LM Studio / Ollama) ---
  # No API key needed. Start your local server first.
  local-lmstudio:
    base_url: http://localhost:1234
    model: lmstudio-community/qwen2.5-coder-7b
    api_key_env: ""
    clear_api_key: true
    tier: sonnet
`
