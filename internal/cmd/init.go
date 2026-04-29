package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project]",
		Short: "Initialize a new envchain vault for a project",
		Long: `Initialize a new envchain vault for the given project name.
If no project name is provided, the current directory name is used.
You will be prompted to set a master password for the vault.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := ""
			if len(args) == 1 {
				project = args[0]
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("could not determine current directory: %w", err)
				}
				project = filepath.Base(cwd)
			}

			if project == "" || project == "." {
				return fmt.Errorf("could not determine project name; please provide one explicitly")
			}

			p := prompter
			password, err := p.AskPassword("Set master password", true)
			if err != nil {
				return fmt.Errorf("password input failed: %w", err)
			}

			v, err := openVaultWithPassword(project, password)
			if err != nil {
				return fmt.Errorf("failed to initialize vault: %w", err)
			}

			if err := v.Save(); err != nil {
				return fmt.Errorf("failed to save vault: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Initialized envchain vault for project %q\n", project)
			return nil
		},
	}
	return cmd
}
