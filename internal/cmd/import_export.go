package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain-cli/internal/vault"
)

func newDumpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dump",
		Short: "Print all vault entries as JSON (plaintext)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := openVault(cmd)
			if err != nil {
				return err
			}

			data, err := v.Export()
			if err != nil {
				return fmt.Errorf("exporting vault: %w", err)
			}

			b, err := vault.MarshalExportData(data)
			if err != nil {
				return fmt.Errorf("serialising export data: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
}

func newLoadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "load <file>",
		Short: "Import entries from a JSON dump file into the vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := openVault(cmd)
			if err != nil {
				return err
			}

			b, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("reading file %q: %w", args[0], err)
			}

			data, err := vault.UnmarshalExportData(b)
			if err != nil {
				return err
			}

			n, err := v.Import(data)
			if err != nil {
				return fmt.Errorf("importing into vault: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Imported %d variable(s).\n", n)
			return nil
		},
	}
}
