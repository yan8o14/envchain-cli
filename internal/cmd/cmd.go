// Package cmd provides the CLI command definitions and routing for envchain.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain-cli/internal/prompt"
	"envchain-cli/internal/vault"
)

// RootCmd is the base command for the envchain CLI.
var RootCmd = &cobra.Command{
	Use:   "envchain",
	Short: "Manage and inject environment variables per project using encrypted local storage",
}

func init() {
	RootCmd.AddCommand(newSetCmd())
	RootCmd.AddCommand(newGetCmd())
	RootCmd.AddCommand(newListCmd())
	RootCmd.AddCommand(newDeleteCmd())
	RootCmd.AddCommand(newExecCmd())
}

// openVault opens the vault for the given namespace, prompting the user for a password.
func openVault(namespace string) (*vault.Vault, error) {
	p := prompt.NewDefault()
	password, err := p.AskPassword("Enter vault password", false)
	if err != nil {
		return nil, fmt.Errorf("reading password: %w", err)
	}
	v, err := vault.New(namespace, password)
	if err != nil {
		return nil, fmt.Errorf("opening vault: %w", err)
	}
	return v, nil
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
