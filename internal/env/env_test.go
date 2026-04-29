package env_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envchain-cli/internal/env"
)

func TestBuildEnvContainsInjectedVars(t *testing.T) {
	vars := map[string]string{
		"MY_SECRET": "supersecret",
		"API_KEY":   "abc123",
	}

	injector := env.New(vars)
	result := injector.BuildEnv()

	found := make(map[string]string)
	for _, e := range result {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			found[parts[0]] = parts[1]
		}
	}

	if found["MY_SECRET"] != "supersecret" {
		t.Errorf("expected MY_SECRET=supersecret, got %s", found["MY_SECRET"])
	}
	if found["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", found["API_KEY"])
	}
}

func TestBuildEnvOverridesExistingVars(t *testing.T) {
	os.Setenv("OVERRIDE_TEST", "original")
	t.Cleanup(func() { os.Unsetenv("OVERRIDE_TEST") })

	vars := map[string]string{"OVERRIDE_TEST": "overridden"}
	injector := env.New(vars)
	result := injector.BuildEnv()

	count := 0
	for _, e := range result {
		if strings.HasPrefix(e, "OVERRIDE_TEST=") {
			count++
			if e != "OVERRIDE_TEST=overridden" {
				t.Errorf("expected overridden value, got %s", e)
			}
		}
	}

	if count != 1 {
		t.Errorf("expected exactly 1 occurrence of OVERRIDE_TEST, got %d", count)
	}
}

func TestBuildEnvPreservesExistingVars(t *testing.T) {
	os.Setenv("EXISTING_VAR", "keep_me")
	t.Cleanup(func() { os.Unsetenv("EXISTING_VAR") })

	vars := map[string]string{"NEW_VAR": "new_value"}
	injector := env.New(vars)
	result := injector.BuildEnv()

	found := false
	for _, e := range result {
		if e == "EXISTING_VAR=keep_me" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected EXISTING_VAR to be preserved in environment")
	}
}

func TestRunReturnsErrorForEmptyArgs(t *testing.T) {
	injector := env.New(map[string]string{})
	err := injector.Run([]string{})
	if err == nil {
		t.Error("expected error for empty args, got nil")
	}
}

func TestRunReturnsErrorForUnknownCommand(t *testing.T) {
	injector := env.New(map[string]string{})
	err := injector.Run([]string{"nonexistent_command_xyz"})
	if err == nil {
		t.Error("expected error for unknown command, got nil")
	}
}
