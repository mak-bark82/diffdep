package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleRenameReport() RenameReport {
	return RenameReport{
		Branch: "feature-x",
		Entries: []RenameEntry{
			{
				OldModule: "github.com/acme/old-sdk",
				NewModule: "github.com/acme/new-sdk",
				Version:   "v2.1.0",
				Note:      "possible rename: old-sdk -> new-sdk",
			},
		},
	}
}

func TestSaveRenameReport_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "rename.json")
	if err := SaveRenameReport(path, sampleRenameReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestLoadRenameReport_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rename.json")
	orig := sampleRenameReport()
	if err := SaveRenameReport(path, orig); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRenameReport(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Branch != orig.Branch {
		t.Errorf("branch: got %q want %q", loaded.Branch, orig.Branch)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Fatalf("entries count: got %d want %d", len(loaded.Entries), len(orig.Entries))
	}
	if loaded.Entries[0].OldModule != orig.Entries[0].OldModule {
		t.Errorf("OldModule mismatch: %s", loaded.Entries[0].OldModule)
	}
}

func TestLoadRenameReport_NotFound(t *testing.T) {
	_, err := LoadRenameReport("/tmp/does-not-exist-rename.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRenameReport_BadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("{invalid"), 0o644)
	_, err := LoadRenameReport(path)
	if err == nil {
		t.Fatal("expected JSON decode error")
	}
}
