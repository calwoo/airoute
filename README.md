# airoute

Switch Claude Code between API endpoints and models with a single command.

```
airoute <route>                    # start claude with that route's env vars
airoute <route> -- <command>       # run any command with the env vars set
eval $(airoute env <route>)        # export vars into your current shell
airoute list                       # show configured routes
airoute init                       # create ~/.airoute/config.yaml
```

## Install

Requires Go 1.21+.

```bash
go install github.com/calwoo/airoute@latest
```

Add `~/go/bin` to your PATH if it isn't already:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

## Usage

### Start Claude with a route

```bash
airoute openrouter-sonnet
# Launches claude with ANTHROPIC_BASE_URL, ANTHROPIC_MODEL,
# and ANTHROPIC_AUTH_TOKEN set from the route config.
```

When you run `airoute <route>`, it replaces the current process with `claude` — no background process left behind, signals pass through cleanly.

### Run any command with a route's env vars

```bash
airoute deepinfra-qwen -- npm run test
# Runs npm test with the route's API endpoint and key injected.
```

Useful for scripts, CI steps, or any tooling that reads `ANTHROPIC_*` env vars.

### Export vars into the current shell

```bash
eval $(airoute env openrouter-sonnet)
claude
```

This sets the route's env vars in your current shell session rather than spawning a subprocess. Combine with a shell function for a shorthand:

```bash
function air() { eval $(airoute env "$1"); }
air deepinfra-qwen
claude
```

### Print env vars (dry run)

```bash
airoute env deepinfra-qwen --no-eval
```

Shows what `buildEnv` would set without executing anything. Handy for debugging.

### List configured routes

```bash
airoute list
```

Displays every route, its endpoint URL, model, key source, and tier.

## Setup

```bash
airoute init          # creates ~/.airoute/config.yaml with examples
# edit the file to set your routes and API key env var names
airoute list          # verify
```

## Config

`~/.airoute/config.yaml` — each route specifies an endpoint, model, and which shell variable holds the API key:

```yaml
routes:
  openrouter-sonnet:
    base_url: https://openrouter.ai/api
    model: anthropic/claude-sonnet-4-6
    api_key_env: OPENROUTER_API_KEY   # name of the var in your shell
    clear_api_key: true               # blank ANTHROPIC_API_KEY for non-Anthropic providers
    tier: sonnet                      # sets ANTHROPIC_DEFAULT_SONNET_MODEL for subagents

  deepinfra-qwen:
    base_url: https://api.deepinfra.com/anthropic
    model: Qwen/Qwen3-30B-A3B
    api_key_env: DEEPINFRA_TOKEN
    clear_api_key: true
    tier: sonnet
```

Set the referenced env vars in your shell (e.g. in `~/.zshrc`):

```bash
export OPENROUTER_API_KEY=sk-or-...
export DEEPINFRA_TOKEN=...
```

