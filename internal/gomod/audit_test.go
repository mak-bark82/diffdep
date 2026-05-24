package gomod

import (
	"os"
	"strings"
	"testing"
)

func sampleDiffForAudit() DiffResult {
	return DiffResult{
		Added:   []Dependency{{Module: "github.com/new/pkg", Version: "v1.0.0"}},
		Removed: []Dependency{{Module: "github.com/old/pkg", Version: "v2.1.0"}},
		Changed: []DependencyChange{
			{Module: "github.com/foo/bar", OldVersion: "v1.2.0", NewVersion: "v2.0.0"},
		},
	}
}

func TestNewAuditReport_EntryCount(t *testing.T) {
	report := NewAuditReport("main", sampleDiffForAudit())
	if len(report.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(report.Entries))
	}
}

func TestNewAuditReport_RiskLevels(t *testing.T) {
	report := NewAuditReport("main", sampleDiffForAudit())
	for _, e := range report.Entries {
		if e.Module == "github.com/foo/bar" && e.Risk != "high" {
			t.Errorf("expected high risk for major change, got %s", e.Risk)
		}
		if e.Module == "github.com/old/pkg" && e.Risk != "medium" {
			t.Errorf("expected medium risk for removed, got %s", e.Risk)
		}
		if e.Module == "github.com/new/pkg" && e.Risk != "low" {
			t.Errorf("expected low risk for added, got %s", e.Risk)
		}
	}
}

func TestFormatAuditReport_ContainsBranch(t *testing.T) {
	report := NewAuditReport("feature-x", sampleDiffForAudit())
	out := FormatAuditReport(report)
	if !strings.Contains(out, "feature-x") {
		t.Error("expected branch name in output")
	}
}

func TestFormatAuditReport_ContainsModules(t *testing.T) {
	report := NewAuditReport("main", sampleDiffForAudit())
	out := FormatAuditReport(report)
	if !strings.Contains(out, "github.com/foo/bar") {
		t.Error("expected module name in output")
	}
	if !strings.Contains(out, "CHANGED") {
		t.Error("expected CHANGED label in output")
	}
}

func TestSaveAndLoadAuditReport(t *testing.T) {
	dir := t.TempDir()
	report := NewAuditReport("main", sampleDiffForAudit())
	if err := SaveAuditReport(dir, report); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAuditReport(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Branch != "main" {
		t.Errorf("expected branch main, got %s", loaded.Branch)
	}
	if len(loaded.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(loaded.Entries))
	}
}

func TestLoadAuditReport_Missing(t *testing.T) {
	_, err := LoadAuditReport("/nonexistent/path")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadAuditReport_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/audit.json"
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadAuditReport(dir)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
