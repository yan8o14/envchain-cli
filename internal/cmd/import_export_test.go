package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDumpCmdRegistered(t *testing.T) {
	out, err := executeCommand("help")
	if err != nil {
		t.Fatalf("help error: %v", err)
	}
	if !strings.Contains(out, "dump") {
		t.Error("expected 'dump' in help output")
	}
	if !strings.Contains(out, "load") {
		t.Error("expected 'load' in help output")
	}
}

func TestDumpCmdRejectsArgs(t *testing.T) {
	_, err := executeCommand("dump", "unexpected")
	if err == nil {
		t.Error("expected error when passing args to dump")
	}
}

func TestLoadCmdRequiresExactlyOneArg(t *testing.T) {
	_, err := executeCommand("load")
	if err == nil {
		t.Error("expected error when no file arg provided")
	}
}

func TestDumpCmdOutputIsValidJSON(t *testing.T) {
	t.Setenv("ENVCHAIN_DIR", t.TempDir())

	_, _ = executeCommand("init", "testproject")
	_, _ = executeCommand("set", "ALPHA=hello")

	out, err := executeCommand("dump")
	if err != nil {
		t.Fatalf("dump error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("dump output is not valid JSON: %v\noutput: %s", err, out)
	}
}

func TestLoadCmdImportsEntries(t *testing.T) {
	t.Setenv("ENVCHAIN_DIR", t.TempDir())
	_, _ = executeCommand("init", "loadtest")

	dumpFile := filepath.Join(t.TempDir(), "dump.json")
	content := `{"project":"loadtest","entries":{"FROM_FILE":"imported_value"}}`
	if err := os.WriteFile(dumpFile, []byte(content), 0600); err != nil {
		t.Fatalf("writing dump file: %v", err)
	}

	out, err := executeCommand("load", dumpFile)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected import count in output, got: %q", out)
	}

	val, err := executeCommand("get", "FROM_FILE")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if !strings.Contains(val, "imported_value") {
		t.Errorf("expected imported_value, got: %q", val)
	}
}

func TestLoadCmdRejectsMissingFile(t *testing.T) {
	t.Setenv("ENVCHAIN_DIR", t.TempDir())
	_, _ = executeCommand("init", "proj")

	_, err := executeCommand("load", "/nonexistent/path/dump.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
