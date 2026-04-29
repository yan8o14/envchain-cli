package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envchain-cli/internal/env"
)

func newExecCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "exec -- COMMAND [ARGS...]",
		Short: "Run a command with vault variables injected into the environment",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := openVault(namespace)
			if err != nil {
				return err
			}

			keys := v.Keys()
			vars := make(map[string]string, len(keys))
			for _, k := range keys {
				val, err := v.Get(k)
				if err != nil {
					return fmt.Errorf("reading key %q: %w", k, err)
				}
				vars[k] = val
			}

			e := env.New(vars)
			if err := e.Run(args); err != nil {
				return fmt.Errorf("exec: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Vault namespace (project name)")
	return cmd
}
