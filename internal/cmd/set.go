package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envchain-cli/internal/prompt"
)

func newSetCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "set KEY",
		Short: "Set an environment variable in the vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			p := prompt.NewDefault()
			value, err := p.AskPassword(fmt.Sprintf("Value for %s", key), false)
			if err != nil {
				return fmt.Errorf("reading value: %w", err)
			}

			v, err := openVault(namespace)
			if err != nil {
				return err
			}

			if err := v.Set(key, value); err != nil {
				return fmt.Errorf("setting key: %w", err)
			}

			fmt.Printf("✓ Set %s in namespace %q\n", key, namespace)
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Vault namespace (project name)")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "delete KEY",
		Short: "Delete an environment variable from the vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			v, err := openVault(namespace)
			if err != nil {
				return err
			}

			if err := v.Delete(key); err != nil {
				return fmt.Errorf("deleting key: %w", err)
			}

			fmt.Printf("✓ Deleted %s from namespace %q\n", key, namespace)
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Vault namespace (project name)")
	return cmd
}
