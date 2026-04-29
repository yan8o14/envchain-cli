package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rename <old-key> <new-key>",
		Short: "Rename an environment variable key",
		Long:  `Rename an existing environment variable key to a new name, preserving its value.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldKey := args[0]
			newKey := args[1]

			if oldKey == newKey {
				return fmt.Errorf("old key and new key must be different")
			}

			v, err := openVault(cmd)
			if err != nil {
				return err
			}

			value, err := v.Get(oldKey)
			if err != nil {
				return fmt.Errorf("key %q not found: %w", oldKey, err)
			}

			if err := v.Set(newKey, value); err != nil {
				return fmt.Errorf("failed to set new key %q: %w", newKey, err)
			}

			if err := v.Delete(oldKey); err != nil {
				return fmt.Errorf("failed to delete old key %q: %w", oldKey, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Renamed %q to %q\n", oldKey, newKey)
			return nil
		},
	}
}
