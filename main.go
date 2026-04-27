package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// knownSubcommands lists all named subcommands so the fallback dispatcher
// can distinguish "airoute list" from "airoute some-route-key".
var knownSubcommands = map[string]bool{
	"init":        true,
	"list":        true,
	"env":         true,
	"help":        true,
	"completion":  true,
}

func main() {
	// If the first arg looks like a route key (not a known subcommand or flag),
	// rewrite os.Args to route through the hidden _run command.
	if len(os.Args) > 1 {
		first := os.Args[1]
		if !knownSubcommands[first] && first != "--help" && first != "-h" {
			os.Args = append([]string{os.Args[0], "_run"}, os.Args[1:]...)
		}
	}

	root := &cobra.Command{
		Use:           "airoute",
		Short:         "Switch Claude Code between API endpoints and models",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(cmdInit(), cmdList(), cmdEnv(), cmdRun())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "airoute: error:", err)
		os.Exit(1)
	}
}

func cmdInit() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create ~/.airoute/config.yaml",
		Long: `Create ~/.airoute/config.yaml with a starter route.

Each route in the config supports the following fields:

  base_url      (required) API endpoint URL
  model         (required) model ID string for the provider
  api_key_env   name of the shell env var holding your API key ("" for none)
  clear_api_key set ANTHROPIC_API_KEY="" — required for non-Anthropic providers
  tier          opus | sonnet | haiku — sets ANTHROPIC_DEFAULT_{TIER}_MODEL

See examples/config.yaml in the repo for provider-specific examples.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInit()
		},
	}
}

func cmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show configured routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doList()
		},
	}
}

func cmdEnv() *cobra.Command {
	return &cobra.Command{
		Use:   "env <route>",
		Short: "Print export statements for a route (for eval usage)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doEnv(args[0])
		},
	}
}

func cmdRun() *cobra.Command {
	return &cobra.Command{
		Use:    "_run <route> [-- command [args...]]",
		Hidden: true, // invoked via the os.Args rewrite, not directly
		Args:   cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			routeKey := args[0]
			passthrough := args[1:]
			// Strip a leading "--" separator if present.
			if len(passthrough) > 0 && passthrough[0] == "--" {
				passthrough = passthrough[1:]
			}
			return doRun(routeKey, passthrough)
		},
	}
}
