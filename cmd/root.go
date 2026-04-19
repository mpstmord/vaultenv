package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultenv",
	aultenv injects HashiCorp Vault secrets into processong: `vaultenv fet them as
enonment variables into a subprocess, keeping secrets out of shell config files.`,
}//unc Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
