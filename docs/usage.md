# airoute

Switch Claude Code between API endpoints and models with a single command. airoute injects the right environment variables so Claude Code connects to whatever provider you want — Anthropic direct, OpenRouter, DeepInfra, Fireworks, or a local model.

## Installation

**Prerequisite:** Go 1.21+.

```bash
go install github.com/calwoo/airoute@latest
```

Make sure `~/go/bin` is on your `PATH`:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

Verify it works:

```bash
airoute --help
```

## Quick Start

```bash
airoute init                  # creates ~/.airoute/config.yaml with a starter route
# edit ~/.airoute/config.yaml to set your routes and API key env var names
airoute list                  # verify routes load
airoute <route-name>          # launch Claude Code with that route
```

## Commands

### `airoute <route>` — Start Claude with a route

Resolves `<route>` from your config, injects its environment variables, and replaces the current process with `claude`. No background process is left behind — signals (Ctrl-C, SIGTERM) pass through directly.

```bash
airoute openrouter-sonnet     # launches claude with OpenRouter's endpoint and key
```

**Run a different command** with the route's env vars by appending `--`:

```bash
airoute deepinfra-qwen -- npm run test
airoute local-ollama -- curl -X POST http://localhost:11434/api/generate
```

The first argument after `--` is the binary to exec; everything after is its arguments.

### `airoute list` — Show configured routes

Prints a table with columns for route name, base URL, model, and tier.

```bash
$ airoute list
ROUTE               BASE_URL                           MODEL                           TIER
anthropic           https://api.anthropic.com           claude-sonnet-4-6                sonnet
deepinfra-qwen      https://api.deepinfra.com/anthropic Qwen/Qwen3-30B-A3B              sonnet
```

### `airoute env <route>` — Print export statements

Prints `export KEY='value'` lines for the route's env vars, designed for `eval`:

```bash
eval $(airoute env openrouter-sonnet)
claude
```

The output is deterministic — keys appear in a fixed order. See [Shell Integration](#shell-integration) for convenience patterns.

### `airoute init` — Create config

Creates `~/.airoute/config.yaml` and the `~/.airoute/` directory with a starter route. If the file already exists, prompts for confirmation before overwriting.

### `airoute help`, `airoute completion`

Standard cobra-generated commands. `help` prints usage; `completion` generates shell completion scripts.

## Configuration

### Config file location

`~/.airoute/config.yaml` — always user-global, never project-local.

### Route fields

Each entry under `routes:` supports the following fields:

| Field | Required | Description |
|---|---|---|
| `base_url` | yes | API endpoint URL (e.g. `https://api.anthropic.com`) |
| `model` | yes | Model ID (e.g. `claude-sonnet-4-6`, `Qwen/Qwen3-30B-A3B`) |
| `api_key_env` | no | Name of the shell env var holding your API key. If empty or omitted, no auth token is set (for local models). |
| `clear_api_key` | no | When `true`, sets `ANTHROPIC_API_KEY=""`. Required for non-Anthropic providers to prevent the SDK from falling back to a real Anthropic key. |
| `tier` | no | One of `opus`, `sonnet`, or `haiku`. Sets the corresponding `ANTHROPIC_DEFAULT_{TIER}_MODEL` env var so Claude Code sub-agents use the same model. |

### Validation rules

- `base_url` and `model` must be non-empty.
- `tier` must be `opus`, `sonnet`, or `haiku`.
- If `api_key_env` references a shell variable that isn't set, airoute prints a clear error pointing you to export it.

### Environment variables injected

| Env var | When set | Value |
|---|---|---|
| `ANTHROPIC_BASE_URL` | always | `route.base_url` |
| `ANTHROPIC_MODEL` | always | `route.model` |
| `ANTHROPIC_AUTH_TOKEN` | `api_key_env` is non-empty | Value of the referenced shell variable |
| `ANTHROPIC_API_KEY` | `clear_api_key: true` | `""` (empty string) |
| `ANTHROPIC_DEFAULT_OPUS_MODEL` | `tier: opus` | `route.model` |
| `ANTHROPIC_DEFAULT_SONNET_MODEL` | `tier: sonnet` | `route.model` |
| `ANTHROPIC_DEFAULT_HAIKU_MODEL` | `tier: haiku` | `route.model` |

Note: `ANTHROPIC_API_KEY` is set to `""` (not unset) because Claude Code falls back to `ANTHROPIC_API_KEY` when `ANTHROPIC_AUTH_TOKEN` is absent. Blanking it forces the SDK to use `ANTHROPIC_AUTH_TOKEN` exclusively.

## Provider Examples

See [examples/config.yaml](../examples/config.yaml) in the repo for ready-to-use route definitions covering these patterns:

| Pattern | `api_key_env` | `clear_api_key` | Notes |
|---|---|---|---|
| **Direct Anthropic** | `ANTHROPIC_API_KEY` | not set | Minimal config, no blanking needed |
| **OpenRouter / DeepInfra / Fireworks** | your API key env var | `true` | Third-party providers with Anthropic-compatible APIs; blanking prevents key fallback |
| **Local (LM Studio, Ollama)** | `""` | `true` | No auth needed; blanking avoids auth token confusion |

The config file generated by `airoute init` is a minimal starting point. The examples file is the expanded reference.

## Shell Integration

### eval shorthand

Load a route's env vars into your current shell session without spawning a child process:

```bash
eval $(airoute env openrouter-sonnet)
claude
```

### Shell function

Add this to your `~/.zshrc` or `~/.bashrc`:

```bash
function air() { eval $(airoute env "$1"); }
```

Then use it from anywhere:

```bash
air deepinfra-qwen
claude
```

### API keys

Set the API key env vars your routes reference, also in your shell rc file:

```bash
export OPENROUTER_API_KEY=sk-or-...
export DEEPINFRA_TOKEN=...
```

## Troubleshooting

| Error | Cause | Fix |
|---|---|---|
| `config not found — run: airoute init` | `~/.airoute/config.yaml` doesn't exist | Run `airoute init` |
| `config must have a 'routes' mapping` | Config file exists but has no `routes:` key | Check YAML structure |
| `route "X" not found. available: ...` | Typo in route name | Run `airoute list` to see valid names |
| `route "X": missing required field: base_url` | Route is missing `base_url` | Add the field to the route |
| `route "X": tier must be opus, sonnet, or haiku` | Invalid tier value | Correct the tier field |
| `'claude' not found in PATH` | Claude Code isn't installed globally | `npm install -g @anthropic-ai/claude-code` |
| `env var "X" is not set in your shell` | Route references an env var that isn't exported | Export it in your shell rc file |
| `config parse error: ...` | Malformed YAML in config file | Validate with a YAML linter or check for syntax errors |