package vault_test

import (
	"os"
	"sort"
	"testing"

	"github.com/envchain-cli/internal/storage"
	"github.com/envchain-cli/internal/vault"
)

func newTestVault(t *testing.T, project, password string) *vault.Vault {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-vault-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	sm := storage.NewStorageManager(dir)
	return vault.New(project, password, sm)
}

func TestSetAndGet(t *testing.T) {
	v := newTestVault(t, "myproject", "supersecret")

	if err := v.Set("DB_URL", "postgres://localhost/mydb"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := v.Get("DB_URL")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "postgres://localhost/mydb" {
		t.Errorf("expected %q, got %q", "postgres://localhost/mydb", val)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	v := newTestVault(t, "myproject", "supersecret")

	_, err := v.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestSetEmptyKeyReturnsError(t *testing.T) {
	v := newTestVault(t, "myproject", "supersecret")

	if err := v.Set("", "value"); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestDeleteKey(t *testing.T) {
	v := newTestVault(t, "myproject", "supersecret")

	if err := v.Set("API_KEY", "abc123"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if err := v.Delete("API_KEY"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := v.Get("API_KEY")
	if err == nil {
		t.Fatal("expected error after deletion, got nil")
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	v := newTestVault(t, "myproject", "supersecret")

	if err := v.Delete("GHOST"); err == nil {
		t.Fatal("expected error deleting non-existent key, got nil")
	}
}

func TestList(t *testing.T) {
	v := newTestVault(t, "myproject", "supersecret")

	keys := []string{"FOO", "BAR", "BAZ"}
	for _, k := range keys {
		if err := v.Set(k, "value"); err != nil {
			t.Fatalf("Set(%q) failed: %v", k, err)
		}
	}

	got, err := v.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	sort.Strings(keys)
	sort.Strings(got)

	if len(got) != len(keys) {
		t.Fatalf("expected %d keys, got %d", len(keys), len(got))
	}
	for i := range keys {
		if keys[i] != got[i] {
			t.Errorf("expected key %q, got %q", keys[i], got[i])
		}
	}
}

func TestWrongPasswordFailsDecrypt(t *testing.T) {
	dir, err := os.MkdirTemp("", "envchain-vault-wrongpw-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	sm := storage.NewStorageManager(dir)
	v1 := vault.New("proj", "correctpassword", sm)
	v2 := vault.New("proj", "wrongpassword", sm)

	if err := v1.Set("SECRET", "topsecret"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	_, err = v2.Get("SECRET")
	if err == nil {
		t.Fatal("expected decryption to fail with wrong password, got nil")
	}
}
