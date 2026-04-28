package vault

import (
	"errors"
	"fmt"

	"github.com/envchain-cli/internal/crypto"
	"github.com/envchain-cli/internal/storage"
)

// Vault manages encrypted environment variables for a named project.
type Vault struct {
	projectName string
	password    string
	sm          *storage.StorageManager
}

// New creates a new Vault for the given project and password.
func New(projectName, password string, sm *storage.StorageManager) *Vault {
	return &Vault{
		projectName: projectName,
		password:    password,
		sm:          sm,
	}
}

// Set encrypts and stores an environment variable for the project.
func (v *Vault) Set(key, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}

	store, err := v.sm.LoadStore(v.projectName)
	if err != nil {
		return fmt.Errorf("loading store: %w", err)
	}

	derived, salt, err := crypto.DeriveKey(v.password, nil)
	if err != nil {
		return fmt.Errorf("deriving key: %w", err)
	}

	encrypted, nonce, err := crypto.Encrypt(derived, []byte(value))
	if err != nil {
		return fmt.Errorf("encrypting value: %w", err)
	}

	store[key] = storage.EncryptedEntry{
		Ciphertext: encrypted,
		Nonce:      nonce,
		Salt:       salt,
	}

	if err := v.sm.SaveStore(v.projectName, store); err != nil {
		return fmt.Errorf("saving store: %w", err)
	}
	return nil
}

// Get decrypts and returns the environment variable value for the given key.
func (v *Vault) Get(key string) (string, error) {
	store, err := v.sm.LoadStore(v.projectName)
	if err != nil {
		return "", fmt.Errorf("loading store: %w", err)
	}

	entry, ok := store[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in project %q", key, v.projectName)
	}

	derived, _, err := crypto.DeriveKey(v.password, entry.Salt)
	if err != nil {
		return "", fmt.Errorf("deriving key: %w", err)
	}

	plaintext, err := crypto.Decrypt(derived, entry.Ciphertext, entry.Nonce)
	if err != nil {
		return "", fmt.Errorf("decrypting value: %w", err)
	}

	return string(plaintext), nil
}

// Delete removes an environment variable from the project store.
func (v *Vault) Delete(key string) error {
	store, err := v.sm.LoadStore(v.projectName)
	if err != nil {
		return fmt.Errorf("loading store: %w", err)
	}

	if _, ok := store[key]; !ok {
		return fmt.Errorf("key %q not found in project %q", key, v.projectName)
	}

	delete(store, key)

	if err := v.sm.SaveStore(v.projectName, store); err != nil {
		return fmt.Errorf("saving store: %w", err)
	}
	return nil
}

// List returns all stored environment variable keys for the project.
func (v *Vault) List() ([]string, error) {
	store, err := v.sm.LoadStore(v.projectName)
	if err != nil {
		return nil, fmt.Errorf("loading store: %w", err)
	}

	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}
	return keys, nil
}
