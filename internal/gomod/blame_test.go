package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForBlame() DiffResult {
	return DiffResult{
		Added: []Dependency{
			{Module: "github.com/new/lib", Version: "v1.0.0"},
		},
		Removed: []Dependency{
			{Module: "github.com/old/lib", Version: "v0.9.0"},
		},
		Changed: []DependencyChange{
			{Module: "github.com/foo/bar", OldVersion: "v1.2.0", NewVersion: "v2.0.0"},
			{Module: "github.com/foo/baz", OldVersion: "v1.5.0", NewVersion: "v1.4.0"},
		},
	}
}

func TestNewBlameReport_EntryCount(t *testing.T) {
	diff := sampleDiffForBlame()
	report := NewBlameReport("feature-x", diff)
	if len(report.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(report.Entries))
	}
}

func TestNewBlameReport_Branch(t *testing.T) {
	report := NewBlameReport("main", sampleDiffForBlame())
	if report.Branch != "main" {
		t.Errorf("expected branch 'main', got %q", report.Branch)
	}
}

func TestNewBlameReport_MajorNote(t *testing.T) {
	diff := sampleDiffForBlame()
	report := NewBlameReport("dev", diff)
	for _, e := range report.Entries {
		if e.Module == "github.com/foo/bar" {
			if e.Note != "major version change" {
				t.Errorf("expected 'major version change', got %q", e.Note)
			}
			return
		}
	}
	t.Error("entry for github.com/foo/bar not found")
}

func TestNewBlameReport_DowngradeNote(t *testing.T) {
	diff := sampleDiffForBlame()
	report := NewBlameReport("dev", diff)
	for _, e := range report.Entries {
		if e.Module == "github.com/foo/baz" {
			if e.Note != "downgrade detected" {
				t.Errorf("expected 'downgrade detected', got %q", e.Note)
			}
			return
		}
	}
	t.Error("entry for github.com/foo/baz not found")
}

func TestFormatBlameReport_ContainsBranch(t *testing.T) {
	report := NewBlameReport("feature-x", sampleDiffForBlame())
	out := FormatBlameReport(report)
	if !strings.Contains(out, "feature-x") {
		t.Errorf("expected branch name in output, got:\n%s", out)
	}
}

func TestFormatBlameReport_Empty(t *testing.T) {
	report := NewBlameReport("main", DiffResult{})
	out := FormatBlameReport(report)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' message, got:\n%s", out)
	}
}

func TestFormatBlameReport_ContainsModules(t *testing.T) {
	report := NewBlameReport("dev", sampleDiffForBlame())
	out := FormatBlameReport(report)
	for _, mod := range []string{"github.com/new/lib", "github.com/old/lib", "github.com/foo/bar"} {
		if !strings.Contains(out, mod) {
			t.Errorf("expected module %q in output", mod)
		}
	}
}
