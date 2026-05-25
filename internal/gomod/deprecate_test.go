package gomod

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleDepsForDeprecate() []Dependency {
	return []Dependency{
		{Module: "github.com/dgrijalva/jwt-go", Version: "v3.2.0+incompatible"},
		{Module: "github.com/some/safe-lib", Version: "v1.0.0"},
		{Module: "github.com/codegangsta/cli", Version: "v1.22.0"},
	}
}

func TestCheckDeprecations_DetectsKnown(t *testing.T) {
	deps := sampleDepsForDeprecate()
	report := CheckDeprecations(deps, "main")
	if len(report.Entries) != 2 {
		t.Fatalf("expected 2 deprecation entries, got %d", len(report.Entries))
	}
}

func TestCheckDeprecations_SkipsSafe(t *testing.T) {
	deps := []Dependency{
		{Module: "github.com/some/safe-lib", Version: "v1.0.0"},
	}
	report := CheckDeprecations(deps, "main")
	if len(report.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(report.Entries))
	}
}

func TestCheckDeprecations_Empty(t *testing.T) {
	report := CheckDeprecations(nil, "feature")
	if len(report.Entries) != 0 {
		t.Errorf("expected empty report, got %d entries", len(report.Entries))
	}
}

func TestFormatDeprecationReport_NoEntries(t *testing.T) {
	r := DeprecationReport{Branch: "main"}
	out := FormatDeprecationReport(r)
	if !strings.Contains(out, "no known deprecated") {
		t.Errorf("expected no-deprecated message, got: %s", out)
	}
}

func TestFormatDeprecationReport_ContainsModule(t *testing.T) {
	r := CheckDeprecations(sampleDepsForDeprecate(), "main")
	out := FormatDeprecationReport(r)
	if !strings.Contains(out, "dgrijalva/jwt-go") {
		t.Errorf("expected module in report, got: %s", out)
	}
}

func TestSaveAndLoadDeprecationReport(t *testing.T) {
	dir := t.TempDir()
	r := CheckDeprecations(sampleDepsForDeprecate(), "main")
	if err := SaveDeprecationReport(dir, r); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := LoadDeprecationReport(dir)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Entries) != len(r.Entries) {
		t.Errorf("entry count mismatch: got %d, want %d", len(loaded.Entries), len(r.Entries))
	}
}

func TestLoadDeprecationReport_Missing(t *testing.T) {
	_, err := LoadDeprecationReport(filepath.Join(os.TempDir(), "nonexistent_deprecate_dir"))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadDeprecationReport_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deprecation_report.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := LoadDeprecationReport(dir)
	if err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}
