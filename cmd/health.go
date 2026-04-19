package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultenv/internal/health"
)

var healthTimeout time.Duration

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check connectivity to Vault and report status",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := os.Getenv("VAULT_ADDR")
		if addr == "" {
			return fmt.Errorf("VAULT_ADDR environment variable is required")
		}

		vc := health.NewVaultChecker(addr)
		runner := health.NewRunner(healthTimeout, vc)
		report := runner.Run(context.Background())

		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return fmt.Errorf("encoding report: %w", err)
		}

		if report.Status != string(health.StatusOK) {
			return fmt.Errorf("health check failed")
		}
		return nil
	},
}

func init() {
	healthCmd.Flags().DurationVar(&healthTimeout, "timeout", 5*time.Second, "per-check timeout")
	rootCmd.AddCommand(healthCmd)
}
