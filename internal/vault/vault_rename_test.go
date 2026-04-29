package vault_test

import (
	"testing"
)

func TestRenameKey(t *testing.T) {
	v := newTestVault(t)

	if err := v.Set("ORIGINAL", "myvalue"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := v.Get("ORIGINAL")
	if err != nil {
		t.Fatalf("Get ORIGINAL failed: %v", err)
	}

	if err := v.Set("RENAMED", val); err != nil {
		t.Fatalf("Set RENAMED failed: %v", err)
	}
	if err := v.Delete("ORIGINAL"); err != nil {
		t.Fatalf("Delete ORIGINAL failed: %v", err)
	}

	got, err := v.Get("RENAMED")
	if err != nil {
		t.Fatalf("Get RENAMED failed: %v", err)
	}
	if got != "myvalue" {
		t.Errorf("expected 'myvalue', got %q", got)
	}

	_, err = v.Get("ORIGINAL")
	if err == nil {
		t.Error("expected error when getting deleted key ORIGINAL")
	}
}

func TestRenamePreservesOtherKeys(t *testing.T) {
	v := newTestVault(t)

	if err := v.Set("KEY_A", "alpha"); err != nil {
		t.Fatalf("Set KEY_A failed: %v", err)
	}
	if err := v.Set("KEY_B", "beta"); err != nil {
		t.Fatalf("Set KEY_B failed: %v", err)
	}

	val, _ := v.Get("KEY_A")
	_ = v.Set("KEY_C", val)
	_ = v.Delete("KEY_A")

	got, err := v.Get("KEY_B")
	if err != nil {
		t.Fatalf("Get KEY_B failed: %v", err)
	}
	if got != "beta" {
		t.Errorf("KEY_B should still be 'beta', got %q", got)
	}

	gotC, err := v.Get("KEY_C")
	if err != nil {
		t.Fatalf("Get KEY_C failed: %v", err)
	}
	if gotC != "alpha" {
		t.Errorf("KEY_C should be 'alpha', got %q", gotC)
	}
}
