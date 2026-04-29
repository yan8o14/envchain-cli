package cmd

import (
	"strings"
	"testing"
)

func TestRenameCmdRegistered(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "rename <old-key> <new-key>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected rename subcommand to be registered")
	}
}

func TestRenameCmdRequiresExactlyTwoArgs(t *testing.T) {
	_, err := executeCommand("rename", "ONLY_ONE")
	if err == nil {
		t.Error("expected error when only one argument is provided")
	}
}

func TestRenameCmdRejectsSameKeys(t *testing.T) {
	_, err := executeCommand("rename", "KEY", "KEY")
	if err == nil {
		t.Error("expected error when old and new keys are identical")
	}
	if err != nil && !strings.Contains(err.Error(), "different") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRenameCmdRenamesKey(t *testing.T) {
	t.Setenv("ENVCHAIN_PROJECT", "rename-test")
	t.Setenv("ENVCHAIN_PASSWORD", "testpassword")

	_, err := executeCommand("init", "rename-test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, err = executeCommand("set", "OLD_KEY=hello")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	out, err := executeCommand("rename", "OLD_KEY", "NEW_KEY")
	if err != nil {
		t.Fatalf("rename failed: %v", err)
	}
	if !strings.Contains(out, "Renamed") {
		t.Errorf("expected output to contain 'Renamed', got: %s", out)
	}

	getOut, err := executeCommand("get", "NEW_KEY")
	if err != nil {
		t.Fatalf("get NEW_KEY failed: %v", err)
	}
	if !strings.Contains(getOut, "hello") {
		t.Errorf("expected value 'hello' under NEW_KEY, got: %s", getOut)
	}

	_, err = executeCommand("get", "OLD_KEY")
	if err == nil {
		t.Error("expected error when getting deleted OLD_KEY")
	}
}

func TestRenameCmdFailsForMissingKey(t *testing.T) {
	t.Setenv("ENVCHAIN_PROJECT", "rename-missing-test")
	t.Setenv("ENVCHAIN_PASSWORD", "testpassword")

	_, err := executeCommand("init", "rename-missing-test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, err = executeCommand("rename", "NONEXISTENT", "NEW_KEY")
	if err == nil {
		t.Error("expected error when renaming nonexistent key")
	}
}
