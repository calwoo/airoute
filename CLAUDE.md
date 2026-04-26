# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`airoute` — CLI tool that injects environment variables to route Claude Code sessions to different API endpoints (Anthropic, OpenRouter, DeepInfra, Fireworks, local models). Part of the `agent-tooling` monorepo.

## Commands

```bash
go build ./...          # compile
go install .            # install binary to $GOPATH/bin (~/.local/go/bin or ~/go/bin)
go mod tidy             # sync go.mod/go.sum after dependency changes
```

Run without installing:
```bash
go run . <args>
```

## Architecture

Single `main` package — all files are in the root directory with `package main`.

| File | Responsibility |
|---|---|
| `main.go` | Cobra root command, subcommand wiring, `os.Args` rewrite for default route-key dispatch |
| `config.go` | `Config`/`Route` structs, `loadConfig()` (YAML), `getRoute()` (validation + lookup) |
| `env.go` | `buildEnv()` — injects Claude Code env vars into a copy of `os.Environ()`; `envOverrides()` for print-only use |
| `commands.go` | `doInit`, `doList`, `doEnv`, `doRun` — one function per CLI command |
| `template.go` | `initTemplate` string constant — the example YAML written by `airoute init` |

### Key design points

- **`os.Args` rewrite**: `main.go` checks if the first arg is a known subcommand before cobra parses. If not, it inserts `_run` so cobra routes the route-key to `cmdRun`. This enables `airoute <route-key>` without a subcommand name.
- **`syscall.Exec`**: `doRun` calls `syscall.Exec` (not `os.StartProcess`) to replace the current process with `claude`. This means no parent process lingers and signals propagate correctly.
- **Env injection**: `buildEnv` rebuilds `os.Environ()` as a `[]string` slice rather than mutating the live environment, because `syscall.Exec` takes an explicit `envv` argument.
- **Config location**: `~/.airoute/config.yaml` — always user-global, never project-local.

### Env vars set per route

| Variable | Set when |
|---|---|
| `ANTHROPIC_BASE_URL` | always |
| `ANTHROPIC_MODEL` | always |
| `ANTHROPIC_AUTH_TOKEN` | `api_key_env` is non-empty |
| `ANTHROPIC_API_KEY` | `clear_api_key: true` (set to `""`) |
| `ANTHROPIC_DEFAULT_{TIER}_MODEL` | `tier` field is set |
