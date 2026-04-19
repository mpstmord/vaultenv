package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultenv/internal/config"
)

var configFile string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Validate a vaultenv configuration file",
	Long:  `Parse and validate a YAML configuration file without executing any process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			return fmt.Errorf("--config flag is required")
		}
		cfg, err := config.LoadFile(configFile)
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "config %q is valid (%d secret mapping(s))\n",
			configFile, len(cfg.Secrets))
		return nil
	},
}

func init() {
	configCmd.Flags().StringVar(&configFile, "config", "", "path to vaultenv YAML config file")
	rootCmd.AddCommand(configCmd)
}
