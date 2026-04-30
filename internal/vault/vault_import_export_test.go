package vault

import (
	"testing"
)

func TestExportReturnsAllEntries(t *testing.T) {
	v := newTestVault(t)
	_ = v.Set("FOO", "bar")
	_ = v.Set("BAZ", "qux")

	data, err := v.Export()
	if err != nil {
		t.Fatalf("Export() error: %v", err)
	}
	if data.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", data.Entries["FOO"])
	}
	if data.Entries["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", data.Entries["BAZ"])
	}
}

func TestExportEmptyVault(t *testing.T) {
	v := newTestVault(t)
	data, err := v.Export()
	if err != nil {
		t.Fatalf("Export() error: %v", err)
	}
	if len(data.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(data.Entries))
	}
}

func TestImportMergesEntries(t *testing.T) {
	v := newTestVault(t)
	_ = v.Set("EXISTING", "value")

	n, err := v.Import(&ExportData{
		Project: "other",
		Entries: map[string]string{"NEW_KEY": "new_val", "EXISTING": "overwritten"},
	})
	if err != nil {
		t.Fatalf("Import() error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 imported, got %d", n)
	}

	val, _ := v.Get("EXISTING")
	if val != "overwritten" {
		t.Errorf("expected EXISTING=overwritten, got %q", val)
	}
	val, _ = v.Get("NEW_KEY")
	if val != "new_val" {
		t.Errorf("expected NEW_KEY=new_val, got %q", val)
	}
}

func TestImportNilDataReturnsError(t *testing.T) {
	v := newTestVault(t)
	_, err := v.Import(nil)
	if err == nil {
		t.Error("expected error for nil import data")
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := &ExportData{
		Project: "myproject",
		Entries: map[string]string{"A": "1", "B": "2"},
	}

	b, err := MarshalExportData(original)
	if err != nil {
		t.Fatalf("MarshalExportData() error: %v", err)
	}

	parsed, err := UnmarshalExportData(b)
	if err != nil {
		t.Fatalf("UnmarshalExportData() error: %v", err)
	}
	if parsed.Project != original.Project {
		t.Errorf("project mismatch: got %q", parsed.Project)
	}
	if parsed.Entries["A"] != "1" || parsed.Entries["B"] != "2" {
		t.Errorf("entries mismatch: %v", parsed.Entries)
	}
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	_, err := UnmarshalExportData([]byte("not-json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
