package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultenv/internal/env"
	"github.com/your-org/vaultenv/internal/process"
	"github.com/your-org/vaultenv/internal/vault"
)

var execCmd = &cobra.Command{
	Use:   "exec -- <command> [args...]",
	Short: "Execute a command with secrets injected into its environment",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mappings, _ := cmd.Flags().GetStringArray("mapping")
		if len(mappings) == 0 {
			return fmt.Errorf("at least one --mapping is required")
		}

		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		resolver := env.NewResolver(client)
		injected := make(map[string]string, len(mappings))

		for _, m := range mappings {
			mapping, err := env.ParseMapping(m)
			if err != nil {
				return fmt.Errorf("invalid mapping %q: %w", m, err)
			}

			value, err := resolver.Resolve(cmd.Context(), mapping)
			if err != nil {
				return fmt.Errorf("resolve %q: %w", m, err)
			}

			injected[mapping.EnvVar] = value
		}

		merged := process.MergeEnv(os.Environ(), injected)
		runner := process.NewRunner(merged)

		return runner.Run(args[0], args[1:])
	},
}

func init() {
	execCmd.Flags().StringArray("mapping", nil,
		"Mapping in the form ENV_VAR=secret/path#field (repeatable)")
	rootCmd.AddCommand(execCmd)
}
