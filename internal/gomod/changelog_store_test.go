package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadChangelog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "changelog.json")

	c := &Changelog{}
	c.AddEntry("main", []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.0.0", NewVersion: "v2.0.0", ChangeType: "changed"},
	})

	if err := SaveChangelog(path, c); err != nil {
		t.Fatalf("SaveChangelog: %v", err)
	}

	loaded, err := LoadChangelog(path)
	if err != nil {
		t.Fatalf("LoadChangelog: %v", err)
	}
	if loaded.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", loaded.Len())
	}
	entry := loaded.Entries[0]
	if entry.Branch != "main" {
		t.Errorf("expected branch main, got %s", entry.Branch)
	}
	if len(entry.Diff) != 1 || entry.Diff[0].Module != "github.com/foo/bar" {
		t.Errorf("unexpected diff contents: %+v", entry.Diff)
	}
}

func TestLoadChangelog_Missing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")

	c, err := LoadChangelog(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if c.Len() != 0 {
		t.Errorf("expected empty changelog, got %d entries", c.Len())
	}
}

func TestLoadChangelog_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadChangelog(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
