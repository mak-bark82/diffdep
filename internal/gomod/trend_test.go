package gomod

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleDiff() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "", NewVersion: "v1.0.0", ChangeType: ChangeAdded},
		{Module: "github.com/baz/qux", OldVersion: "v1.0.0", NewVersion: "", ChangeType: ChangeRemoved},
		{Module: "github.com/acme/lib", OldVersion: "v1.2.0", NewVersion: "v2.0.0", ChangeType: ChangeUpdated},
	}
}

func TestRecordSnapshot_CountsCorrectly(t *testing.T) {
	var trend Trend
	trend.RecordSnapshot("main", sampleDiff())

	if len(trend.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(trend.Entries))
	}
	e := trend.Entries[0]
	if e.Added != 1 || e.Removed != 1 || e.Changed != 1 {
		t.Errorf("unexpected counts: added=%d removed=%d changed=%d", e.Added, e.Removed, e.Changed)
	}
	if e.Branch != "main" {
		t.Errorf("expected branch main, got %s", e.Branch)
	}
}

func TestTrend_Summary_Empty(t *testing.T) {
	var trend Trend
	got := trend.Summary()
	if got != "no trend data recorded" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestTrend_Summary_ContainsBranch(t *testing.T) {
	var trend Trend
	trend.RecordSnapshot("feature-x", sampleDiff())
	sum := trend.Summary()
	if !strings.Contains(sum, "feature-x") {
		t.Errorf("summary missing branch name: %s", sum)
	}
}

func TestSaveAndLoadTrend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trend.json")

	var trend Trend
	trend.RecordSnapshot("main", sampleDiff())

	if err := SaveTrend(path, &trend); err != nil {
		t.Fatalf("SaveTrend: %v", err)
	}

	loaded, err := LoadTrend(path)
	if err != nil {
		t.Fatalf("LoadTrend: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry after load, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Branch != "main" {
		t.Errorf("branch mismatch: %s", loaded.Entries[0].Branch)
	}
}

func TestLoadTrend_Missing(t *testing.T) {
	loaded, err := LoadTrend("/nonexistent/path/trend.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(loaded.Entries) != 0 {
		t.Errorf("expected empty trend, got %d entries", len(loaded.Entries))
	}
}

func TestLoadTrend_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := LoadTrend(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
