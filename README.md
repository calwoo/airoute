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
git clone https://github.com/calvinwoo/airoute
cd airoute
go install .
```

Add `~/go/bin` to your PATH if it isn't already:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

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

## Shell alias (optional)

For an eval-based workflow that sets vars in your current shell:

```bash
# ~/.zshrc
function air() { eval $(airoute env "$1"); }

# Usage
air deepinfra-qwen
claude
```
