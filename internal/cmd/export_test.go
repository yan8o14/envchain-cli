package cmd

import (
	"strings"
	"testing"
)

func TestExportCmdRegistered(t *testing.T) {
	root := newRootCmdForTest()
	found := false
	for _, sub := range root.Commands() {
		if sub.Name() == "export" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'export' subcommand to be registered")
	}
}

func TestExportCmdRejectsArgs(t *testing.T) {
	out, err := executeCommand("export", "unexpected-arg")
	if err == nil {
		t.Fatalf("expected error for unexpected arg, got output: %s", out)
	}
}

func TestExportCmdPosixOutput(t *testing.T) {
	t.Setenv("ENVCHAIN_PROJECT", "test-export-posix")

	// pre-populate via set
	_, err := executeCommand("set", "MY_TOKEN=abc123")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	out, err := executeCommand("export", "--shell", "posix")
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(out, "export MY_TOKEN=") {
		t.Errorf("expected posix export statement, got: %s", out)
	}
}

func TestExportCmdFishOutput(t *testing.T) {
	t.Setenv("ENVCHAIN_PROJECT", "test-export-fish")

	_, err := executeCommand("set", "FISH_VAR=hello")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	out, err := executeCommand("export", "--shell", "fish")
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(out, "set -x FISH_VAR") {
		t.Errorf("expected fish set statement, got: %s", out)
	}
}

func TestExportCmdEmptyVaultPrintsNotice(t *testing.T) {
	t.Setenv("ENVCHAIN_PROJECT", "test-export-empty")

	// do not set any vars — vault should be empty
	_, err := executeCommand("export")
	// empty vault is not an error
	if err != nil {
		t.Fatalf("unexpected error for empty vault: %v", err)
	}
}
