# Future Improvements

Ideas for evolving airoute, grouped by type and rough priority.

---

## Infrastructure

- **Unit tests for core logic** — `buildEnv`, `getRoute`, `loadConfig`, and `doList` formatting are all pure functions with no test coverage. Top priority.
- **Integration test for exec path** — `syscall.Exec` is hard to test directly, but verifying the assembled environment before exec would catch regressions in `buildEnv`.
- **GitHub Actions CI** — `go build`, `go vet`, and `go test ./...` on every push/PR. ~20 lines of YAML.
- **Add golangci-lint config** — Catches common bugs and style issues, especially valuable while test coverage is low.
- **Dependabot** — Auto-create PRs for dependency updates to `go.mod`, keeping cobra and yaml.v3 current.

## Features (new capabilities)

- **Interactive picker on bare `airoute`** — Instead of dumping help text, show an interactive route picker (numbered or fzf-style). First-time users with no config still see help. Turns a dead end into a launchpad.
- **`airoute add` / `airoute rm`** — Add or remove routes without editing YAML directly. `airoute add my-route --base-url https://... --model claude-sonnet-4-6 --key-env MY_KEY`.
- **JSON/YAML output for `airoute list`** — `airoute list -o json | jq` enables programmatic use in scripts and fzf integration.
- **goreleaser for binary releases** — Prebuilt binaries make installation easier for non-Go users and enable Homebrew taps.

## Developer Experience (polish on existing capabilities)

- **Prefix matching for route names** — `airoute deep` matches `deepinfra-qwen` if unique. Saves typing on long route names.
- **Default route** — A `default` field in config, or auto-select when only one route exists. `airoute` with no args launches it directly.
- **Confirmation splash on launch** — Before `syscall.Exec`, print `→ Starting Claude with "openrouter-sonnet" (OpenRouter / claude-sonnet-4-6)` so you visually confirm the route.
- **`--version` flag** — Cobra supports `root.Version` built-in. Useful for debugging and issue reports.
- **Shell completion install docs** — Cobra generates completion scripts; the docs could show how to install them permanently.
- **Color in `airoute list`** — Green for routes whose key is present, red for missing, gray for local-only (no key needed).
- **Route group tags** — Tag routes (`tag: local`, `tag: cheap`) and filter: `airoute list --tag local`.

## Robustness (preventing bugs)

- **Use `os.UserHomeDir()`** instead of `os.Getenv("HOME")` in `config.go`. Handles cross-platform edge cases (Windows, `HOME` unset).
- **Validate `base_url` format** — `getRoute` checks for empty string but doesn't verify the URL is well-formed. A malformed URL fails silently until Claude can't connect.

## Code Quality (internal hygiene)

- **Derive `knownSubcommands` from Cobra** — `main.go` maintains a manual map that duplicates Cobra's internal command list. Could iterate `root.Commands()` instead.
