package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForGroup() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.0.0", NewVersion: "v1.1.0", ChangeType: ChangeUpdated},
		{Module: "github.com/foo/baz", OldVersion: "", NewVersion: "v0.2.0", ChangeType: ChangeAdded},
		{Module: "github.com/other/lib", OldVersion: "v3.0.0", NewVersion: "", ChangeType: ChangeRemoved},
		{Module: "golang.org/x/text", OldVersion: "v0.3.0", NewVersion: "v0.4.0", ChangeType: ChangeUpdated},
	}
}

func TestGroupByOrg_GroupsCorrectly(t *testing.T) {
	diff := sampleDiffForGroup()
	groups := GroupByOrg(diff)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}

	// Groups should be sorted alphabetically
	if groups[0].Name != "github.com/foo" {
		t.Errorf("expected first group 'github.com/foo', got %q", groups[0].Name)
	}
	if len(groups[0].Entries) != 2 {
		t.Errorf("expected 2 entries in github.com/foo, got %d", len(groups[0].Entries))
	}
}

func TestGroupByOrg_Empty(t *testing.T) {
	groups := GroupByOrg([]DiffEntry{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty diff, got %d", len(groups))
	}
}

func TestGroupByOrg_SingleSegmentModule(t *testing.T) {
	diff := []DiffEntry{
		{Module: "stdlib", OldVersion: "v1.0.0", NewVersion: "v2.0.0", ChangeType: ChangeUpdated},
	}
	groups := GroupByOrg(diff)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "stdlib" {
		t.Errorf("expected group name 'stdlib', got %q", groups[0].Name)
	}
}

func TestFormatGroupReport_ContainsOrg(t *testing.T) {
	diff := sampleDiffForGroup()
	groups := GroupByOrg(diff)
	report := FormatGroupReport(groups)

	if !strings.Contains(report, "github.com/foo") {
		t.Errorf("expected report to contain 'github.com/foo'")
	}
	if !strings.Contains(report, "github.com/foo/bar") {
		t.Errorf("expected report to contain module name")
	}
}

func TestFormatGroupReport_NoChanges(t *testing.T) {
	report := FormatGroupReport([]DiffGroup{})
	if !strings.Contains(report, "No dependency changes") {
		t.Errorf("expected no-changes message, got: %s", report)
	}
}

func TestFormatGroupReport_ChangeSymbols(t *testing.T) {
	diff := sampleDiffForGroup()
	groups := GroupByOrg(diff)
	report := FormatGroupReport(groups)

	if !strings.Contains(report, "+") {
		t.Errorf("expected '+' symbol for added entry")
	}
	if !strings.Contains(report, "-") {
		t.Errorf("expected '-' symbol for removed entry")
	}
	if !strings.Contains(report, "~") {
		t.Errorf("expected '~' symbol for updated entry")
	}
}
