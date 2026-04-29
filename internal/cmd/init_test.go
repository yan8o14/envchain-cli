package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCmdRegistered(t *testing.T) {
	args := []string{"--help"}
	out, err := executeCommand(args...)
	require.NoError(t, err)
	assert.Contains(t, out, "init")
}

func TestInitCmdRejectsMultipleArgs(t *testing.T) {
	_, err := executeCommand("init", "proj1", "proj2")
	assert.Error(t, err)
}

func TestInitCmdUsesCurrentDirName(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "myproject")
	require.NoError(t, os.MkdirAll(projectDir, 0o755))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { os.Chdir(origDir) })

	require.NoError(t, os.Chdir(projectDir))

	var buf bytes.Buffer
	rootCmd := buildTestRoot(tmpDir, "testpass", &buf)
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"init"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "myproject")
}

func TestInitCmdWithExplicitProjectName(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	rootCmd := buildTestRoot(tmpDir, "testpass", &buf)
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"init", "myapp"})

	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "myapp")
}

func TestInitCmdCreatesVaultFile(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	rootCmd := buildTestRoot(tmpDir, "testpass", &buf)
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"init", "newproject"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	vaultPath := filepath.Join(tmpDir, "newproject")
	_, statErr := os.Stat(vaultPath)
	assert.NoError(t, statErr, "vault file should exist after init")
}
