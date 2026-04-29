package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "get KEY",
		Short: "Retrieve an environment variable from the vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			v, err := openVault(namespace)
			if err != nil {
				return err
			}

			value, err := v.Get(key)
			if err != nil {
				return fmt.Errorf("getting key %q: %w", key, err)
			}

			fmt.Println(value)
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Vault namespace (project name)")
	return cmd
}

func newListCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys stored in the vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := openVault(namespace)
			if err != nil {
				return err
			}

			keys := v.Keys()
			if len(keys) == 0 {
				fmt.Printf("No keys found in namespace %q\n", namespace)
				return nil
			}

			fmt.Printf("Keys in namespace %q:\n", namespace)
			for _, k := range keys {
				fmt.Printf("  - %s\n", k)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Vault namespace (project name)")
	return cmd
}
