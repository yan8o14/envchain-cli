package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain-cli/internal/env"
)

func newExportCmd() *cobra.Command {
	var shell string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Print export statements for all variables in the vault",
		Long: `Print shell export statements for all variables stored in the vault.

Usage:
  eval $(envchain export)            # bash/zsh
  envchain export --shell fish | source  # fish`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			vault, err := openVault(cmd)
			if err != nil {
				return err
			}

			keys, err := vault.List()
			if err != nil {
				return fmt.Errorf("listing vault keys: %w", err)
			}

			if len(keys) == 0 {
				fmt.Fprintln(cmd.ErrOrStderr(), "no variables stored in vault")
				return nil
			}

			vars := make(map[string]string, len(keys))
			for _, k := range keys {
				v, err := vault.Get(k)
				if err != nil {
					return fmt.Errorf("reading key %q: %w", k, err)
				}
				vars[k] = v
			}

			return env.Export(os.Stdout, vars, shell)
		},
	}

	cmd.Flags().StringVar(&shell, "shell", "posix", `target shell format: posix, fish`)
	return cmd
}
