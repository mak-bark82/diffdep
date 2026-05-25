package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleDeprecationReport() DeprecationReport {
	return DeprecationReport{
		Branch: "main",
		Entries: []DeprecationEntry{
			{
				Module:   "github.com/dgrijalva/jwt-go",
				Version:  "v3.2.0+incompatible",
				Reason:   "replaced by github.com/golang-jwt/jwt",
				Severity: "warn",
			},
		},
	}
}

func TestSaveDeprecationReport_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	r := sampleDeprecationReport()
	if err := SaveDeprecationReport(dir, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path := filepath.Join(dir, "deprecation_report.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file at %s, not found", path)
	}
}

func TestLoadDeprecationReport_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	orig := sampleDeprecationReport()
	if err := SaveDeprecationReport(dir, orig); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadDeprecationReport(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Branch != orig.Branch {
		t.Errorf("branch mismatch: got %q want %q", loaded.Branch, orig.Branch)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("entry count: got %d want %d", len(loaded.Entries), len(orig.Entries))
	}
	if loaded.Entries[0].Module != orig.Entries[0].Module {
		t.Errorf("module mismatch: got %q want %q", loaded.Entries[0].Module, orig.Entries[0].Module)
	}
}

func TestLoadDeprecationReport_NotFound(t *testing.T) {
	_, err := LoadDeprecationReport(filepath.Join(os.TempDir(), "no_such_deprecate_xyz"))
	if err == nil {
		t.Error("expected error for missing report")
	}
}

func TestLoadDeprecationReport_BadJSON(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "deprecation_report.json"), []byte("{"), 0o644)
	_, err := LoadDeprecationReport(dir)
	if err == nil {
		t.Error("expected JSON parse error")
	}
}
