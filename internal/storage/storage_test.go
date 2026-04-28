package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestManager(t *testing.T) *StorageManager {
	t.Helper()
	tmpDir := t.TempDir()
	return &StorageManager{BaseDir: tmpDir}
}

func TestLoadStoreNotExist(t *testing.T) {
	sm := newTestManager(t)
	store, err := sm.LoadStore("myproject")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Project != "myproject" {
		t.Errorf("expected project 'myproject', got '%s'", store.Project)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty entries, got %v", store.Entries)
	}
}

func TestSaveAndLoadStore(t *testing.T) {
	sm := newTestManager(t)
	store := &Store{
		Project: "testproject",
		Entries: map[string]string{
			"API_KEY": "encryptedvalue123",
			"DB_PASS": "encryptedpass456",
		},
	}
	if err := sm.SaveStore(store); err != nil {
		t.Fatalf("SaveStore failed: %v", err)
	}
	loaded, err := sm.LoadStore("testproject")
	if err != nil {
		t.Fatalf("LoadStore failed: %v", err)
	}
	if loaded.Project != store.Project {
		t.Errorf("project mismatch: got %s", loaded.Project)
	}
	if loaded.Entries["API_KEY"] != "encryptedvalue123" {
		t.Errorf("API_KEY mismatch: got %s", loaded.Entries["API_KEY"])
	}
	if loaded.Entries["DB_PASS"] != "encryptedpass456" {
		t.Errorf("DB_PASS mismatch: got %s", loaded.Entries["DB_PASS"])
	}
}

func TestSaveStoreCreatesDirectoryWithCorrectPermissions(t *testing.T) {
	sm := newTestManager(t)
	store := &Store{Project: "permtest", Entries: map[string]string{"X": "y"}}
	if err := sm.SaveStore(store); err != nil {
		t.Fatalf("SaveStore failed: %v", err)
	}
	info, err := os.Stat(filepath.Join(sm.BaseDir, "permtest", storeFileName))
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file perm 0600, got %v", info.Mode().Perm())
	}
}

func TestDeleteStore(t *testing.T) {
	sm := newTestManager(t)
	store := &Store{Project: "todelete", Entries: map[string]string{"K": "v"}}
	_ = sm.SaveStore(store)
	if err := sm.DeleteStore("todelete"); err != nil {
		t.Fatalf("DeleteStore failed: %v", err)
	}
	_, err := os.Stat(filepath.Join(sm.BaseDir, "todelete"))
	if !os.IsNotExist(err) {
		t.Errorf("expected directory to be removed")
	}
}
