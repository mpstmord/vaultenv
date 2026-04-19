package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"

	"vaultenv/internal/env"
	"vaultenv/internal/vault"
)

var (
	mappings []string
	vaultAddr string
	vaultToken string
)

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Inject secrets from Vault into a process environment",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		r := env.NewResolver(client)

		parsed := make([]env.Mapping, 0, len(mappings))
		for _, m := range mappings {
			mp, err := env.ParseMapping(m)
			if err != nil {
				return fmt.Errorf("parse mapping %q: %w", m, err)
			}
			parsed = append(parsed, mp)
		}

		injected, err := env.InjectIntoEnv(os.Environ(), parsed, r)
		if err != nil {
			return fmt.Errorf("inject secrets: %w", err)
		}

		path, err := exec.LookPath(args[0])
		if err != nil {
			return fmt.Errorf("lookup binary %q: %w", args[0], err)
		}

		return syscall.Exec(path, args, injected)
	},
}

func init() {
	runCmd.Flags().StringArrayVarP(&mappings, "mapping", "m", nil,
		"Env var mapping in format ENV_VAR=secret/path#field (repeatable)")
	runCmd.Flags().StringVar(&vaultAddr, "vault-addr", os.Getenv("VAULT_ADDR"),
		"Vault server address (default: $VAULT_ADDR)")
	runCmd.Flags().StringVar(&vaultToken, "vault-token", os.Getenv("VAULT_TOKEN"),
		"Vault token (default: $VAULT_TOKEN)")
	_ = runCmd.MarkFlagRequired("mapping")
}
