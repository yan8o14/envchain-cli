package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultStorageDir = ".envchain"
const storeFileName = "store.json"

// Store represents the encrypted storage structure for a project.
type Store struct {
	Project string            `json:"project"`
	Entries map[string]string `json:"entries"` // key -> encrypted value (base64)
}

// StorageManager handles reading and writing stores to disk.
type StorageManager struct {
	BaseDir string
}

// NewStorageManager creates a StorageManager rooted at the user's home directory.
func NewStorageManager() (*StorageManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &StorageManager{BaseDir: filepath.Join(home, defaultStorageDir)}, nil
}

// storePath returns the path to the store file for the given project.
func (sm *StorageManager) storePath(project string) string {
	return filepath.Join(sm.BaseDir, project, storeFileName)
}

// LoadStore reads and deserializes the store for the given project.
// Returns an empty store if the file does not exist.
func (sm *StorageManager) LoadStore(project string) (*Store, error) {
	path := sm.storePath(project)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{Project: project, Entries: make(map[string]string)}, nil
	}
	if err != nil {
		return nil, err
	}
	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

// SaveStore serializes and writes the store for the given project to disk.
func (sm *StorageManager) SaveStore(store *Store) error {
	dir := filepath.Join(sm.BaseDir, store.Project)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sm.storePath(store.Project), data, 0600)
}

// DeleteStore removes the store file and directory for the given project.
func (sm *StorageManager) DeleteStore(project string) error {
	dir := filepath.Join(sm.BaseDir, project)
	return os.RemoveAll(dir)
}
