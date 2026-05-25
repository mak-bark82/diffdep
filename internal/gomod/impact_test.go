package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForImpact() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "", NewVersion: "v1.0.0", Kind: KindAdded},
		{Module: "github.com/foo/baz", OldVersion: "v1.2.0", NewVersion: "", Kind: KindRemoved},
		{Module: "github.com/foo/qux", OldVersion: "v1.0.0", NewVersion: "v2.0.0", Kind: KindChanged},
		{Module: "github.com/foo/minor", OldVersion: "v1.0.0", NewVersion: "v1.1.0", Kind: KindChanged},
		{Module: "github.com/foo/down", OldVersion: "v1.5.0", NewVersion: "v1.3.0", Kind: KindChanged},
	}
}

func TestAssessImpact_EntryCount(t *testing.T) {
	diff := sampleDiffForImpact()
	report := AssessImpact("main", diff)
	if len(report.Entries) != len(diff) {
		t.Fatalf("expected %d entries, got %d", len(diff), len(report.Entries))
	}
}

func TestAssessImpact_AddedIsLow(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/bar", NewVersion: "v1.0.0", Kind: KindAdded},
	}
	report := AssessImpact("main", diff)
	if report.Entries[0].Level != ImpactLow {
		t.Errorf("expected low impact for added dep, got %s", report.Entries[0].Level)
	}
}

func TestAssessImpact_MajorBumpIsCritical(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/qux", OldVersion: "v1.0.0", NewVersion: "v2.0.0", Kind: KindChanged},
	}
	report := AssessImpact("feature", diff)
	if report.Entries[0].Level != ImpactCritical {
		t.Errorf("expected critical for major bump, got %s", report.Entries[0].Level)
	}
}

func TestAssessImpact_DowngradeIsHigh(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/down", OldVersion: "v1.5.0", NewVersion: "v1.3.0", Kind: KindChanged},
	}
	report := AssessImpact("main", diff)
	if report.Entries[0].Level != ImpactHigh {
		t.Errorf("expected high impact for downgrade, got %s", report.Entries[0].Level)
	}
}

func TestAssessImpact_RemovedIsMedium(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/baz", OldVersion: "v1.2.0", Kind: KindRemoved},
	}
	report := AssessImpact("main", diff)
	if report.Entries[0].Level != ImpactMedium {
		t.Errorf("expected medium impact for removed dep, got %s", report.Entries[0].Level)
	}
}

func TestFormatImpactReport_Empty(t *testing.T) {
	report := AssessImpact("main", nil)
	out := FormatImpactReport(report)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' in empty report, got: %s", out)
	}
}

func TestFormatImpactReport_ContainsBranch(t *testing.T) {
	report := AssessImpact("release", sampleDiffForImpact())
	out := FormatImpactReport(report)
	if !strings.Contains(out, "release") {
		t.Errorf("expected branch name in report output")
	}
}

func TestFormatImpactReport_ContainsModule(t *testing.T) {
	report := AssessImpact("main", sampleDiffForImpact())
	out := FormatImpactReport(report)
	if !strings.Contains(out, "github.com/foo/qux") {
		t.Errorf("expected module name in report output")
	}
}
