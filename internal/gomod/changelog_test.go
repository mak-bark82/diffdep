package gomod

import (
	"strings"
	"testing"
)

func TestChangelog_AddEntry_Empty(t *testing.T) {
	c := &Changelog{}
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", c.Len())
	}
	c.AddEntry("main", []DiffEntry{})
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", c.Len())
	}
}

func TestChangelog_AddEntry_SetsFields(t *testing.T) {
	c := &Changelog{}
	diff := []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.0.0", NewVersion: "v2.0.0", ChangeType: "changed"},
	}
	c.AddEntry("feature-x", diff)
	entry := c.Entries[0]
	if entry.Branch != "feature-x" {
		t.Errorf("expected branch feature-x, got %s", entry.Branch)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if len(entry.Diff) != 1 {
		t.Errorf("expected 1 diff entry, got %d", len(entry.Diff))
	}
}

func TestChangelog_FormatText_ContainsModule(t *testing.T) {
	c := &Changelog{}
	diff := []DiffEntry{
		{Module: "github.com/pkg/errors", OldVersion: "v0.8.0", NewVersion: "v0.9.0", ChangeType: "changed"},
	}
	c.AddEntry("main", diff)
	out := c.FormatText()
	if !strings.Contains(out, "github.com/pkg/errors") {
		t.Errorf("expected module in output, got:\n%s", out)
	}
	if !strings.Contains(out, "v0.8.0") {
		t.Errorf("expected old version in output, got:\n%s", out)
	}
	if !strings.Contains(out, "v0.9.0") {
		t.Errorf("expected new version in output, got:\n%s", out)
	}
}

func TestChangelog_FormatText_MultipleEntries(t *testing.T) {
	c := &Changelog{}
	c.AddEntry("main", []DiffEntry{
		{Module: "github.com/a/b", OldVersion: "v1.0.0", NewVersion: "v1.1.0", ChangeType: "changed"},
	})
	c.AddEntry("dev", []DiffEntry{
		{Module: "github.com/c/d", OldVersion: "", NewVersion: "v2.0.0", ChangeType: "added"},
	})
	out := c.FormatText()
	if !strings.Contains(out, "branch=main") {
		t.Errorf("expected branch=main in output")
	}
	if !strings.Contains(out, "branch=dev") {
		t.Errorf("expected branch=dev in output")
	}
}
