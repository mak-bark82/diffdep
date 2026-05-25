package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForCompat() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", Type: Changed, OldVersion: "v1.3.0", NewVersion: "v2.0.0"},
		{Module: "github.com/foo/bar/v2", Type: Changed, OldVersion: "v2.1.0", NewVersion: "v3.0.0"},
		{Module: "github.com/baz/qux", Type: Changed, OldVersion: "v1.0.0", NewVersion: "v1.5.0"},
		{Module: "github.com/safe/pkg", Type: Added, OldVersion: "", NewVersion: "v2.0.0"},
	}
}

func TestCheckCompat_DetectsMissingSuffix(t *testing.T) {
	diff := sampleDiffForCompat()
	report := CheckCompat("main", diff)
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 compat entry, got %d", len(report.Entries))
	}
	if report.Entries[0].Module != "github.com/foo/bar" {
		t.Errorf("unexpected module: %s", report.Entries[0].Module)
	}
}

func TestCheckCompat_IgnoresCorrectSuffix(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/bar/v2", Type: Changed, OldVersion: "v2.0.0", NewVersion: "v3.0.0"},
	}
	report := CheckCompat("main", diff)
	// bar/v2 going to v3 should still flag (suffix says v2, version is v3)
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
}

func TestCheckCompat_IgnoresMinorBump(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/baz/qux", Type: Changed, OldVersion: "v1.0.0", NewVersion: "v1.5.0"},
	}
	report := CheckCompat("feature", diff)
	if len(report.Entries) != 0 {
		t.Errorf("expected no entries for minor bump, got %d", len(report.Entries))
	}
}

func TestCheckCompat_IgnoresAdded(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/new/mod", Type: Added, OldVersion: "", NewVersion: "v2.0.0"},
	}
	report := CheckCompat("main", diff)
	if len(report.Entries) != 0 {
		t.Errorf("expected no entries for added deps, got %d", len(report.Entries))
	}
}

func TestFormatCompatReport_NoIssues(t *testing.T) {
	report := CompatReport{Branch: "main", Entries: nil}
	out := FormatCompatReport(report)
	if !strings.Contains(out, "no issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestFormatCompatReport_ContainsModule(t *testing.T) {
	report := CompatReport{
		Branch: "main",
		Entries: []CompatEntry{
			{Module: "github.com/foo/bar", From: "v1.0.0", To: "v2.0.0", GoMajor: true, Reason: "major version bump (v1.0.0 -> v2.0.0) without import path suffix"},
		},
	}
	out := FormatCompatReport(report)
	if !strings.Contains(out, "github.com/foo/bar") {
		t.Errorf("expected module name in output, got: %s", out)
	}
	if !strings.Contains(out, "COMPAT") {
		t.Errorf("expected COMPAT tag in output, got: %s", out)
	}
}
