package vault

import (
	"encoding/json"
	"fmt"
)

// ExportData represents a portable, unencrypted snapshot of vault entries.
type ExportData struct {
	Project string            `json:"project"`
	Entries map[string]string `json:"entries"`
}

// Export returns an ExportData snapshot of all key/value pairs in the vault.
// The caller is responsible for handling the sensitivity of the returned data.
func (v *Vault) Export() (*ExportData, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	entries := make(map[string]string, len(v.store.Entries))
	for k, val := range v.store.Entries {
		entries[k] = val
	}

	return &ExportData{
		Project: v.store.Project,
		Entries: entries,
	}, nil
}

// Import merges entries from the given ExportData into the vault.
// Existing keys are overwritten. The vault is saved after a successful import.
func (v *Vault) Import(data *ExportData) (imported int, err error) {
	if data == nil {
		return 0, fmt.Errorf("import data must not be nil")
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	for k, val := range data.Entries {
		if k == "" {
			continue
		}
		v.store.Entries[k] = val
		imported++
	}

	if err := v.save(); err != nil {
		return 0, fmt.Errorf("saving vault after import: %w", err)
	}

	return imported, nil
}

// MarshalExportData serialises ExportData to JSON bytes.
func MarshalExportData(data *ExportData) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// UnmarshalExportData deserialises JSON bytes into ExportData.
func UnmarshalExportData(b []byte) (*ExportData, error) {
	var data ExportData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("parsing export data: %w", err)
	}
	return &data, nil
}
