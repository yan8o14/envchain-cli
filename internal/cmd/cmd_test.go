package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand is a test helper that runs a cobra command with the given args
// and returns the combined output and any error.
func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return buf.String(), err
}

func TestRootCmdHasExpectedSubcommands(t *testing.T) {
	subcmds := map[string]bool{}
	for _, c := range RootCmd.Commands() {
		subcmds[c.Name()] = true
	}

	expected := []string{"set", "get", "list", "delete", "exec"}
	for _, name := range expected {
		if !subcmds[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestSetCmdRequiresExactlyOneArg(t *testing.T) {
	cmd := newSetCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when no args provided to set, got nil")
	}
}

func TestGetCmdRequiresExactlyOneArg(t *testing.T) {
	cmd := newGetCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when no args provided to get, got nil")
	}
}

func TestDeleteCmdRequiresExactlyOneArg(t *testing.T) {
	cmd := newDeleteCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when no args provided to delete, got nil")
	}
}

func TestExecCmdRequiresAtLeastOneArg(t *testing.T) {
	cmd := newExecCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when no args provided to exec, got nil")
	}
}

func TestListCmdAcceptsNoArgs(t *testing.T) {
	cmd := newListCmd()
	if cmd.Args == nil {
		t.Error("expected list cmd to have Args validator")
	}
}
