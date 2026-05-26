package gomod

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleRenameBase() map[string]string {
	return map[string]string{
		"github.com/old/logger": "v1.2.0",
		"github.com/stable/pkg": "v2.0.0",
	}
}

func sampleRenameHead() map[string]string {
	return map[string]string{
		"github.com/new/logger": "v1.2.0",
		"github.com/stable/pkg": "v2.0.0",
	}
}

func TestDetectRenames_FindsRename(t *testing.T) {
	report := DetectRenames(sampleRenameBase(), sampleRenameHead(), "main")
	if len(report.Entries) == 0 {
		t.Fatal("expected at least one rename entry")
	}
	e := report.Entries[0]
	if e.OldModule != "github.com/old/logger" {
		t.Errorf("unexpected OldModule: %s", e.OldModule)
	}
	if e.NewModule != "github.com/new/logger" {
		t.Errorf("unexpected NewModule: %s", e.NewModule)
	}
}

func TestDetectRenames_NoRename(t *testing.T) {
	base := map[string]string{"github.com/foo/bar": "v1.0.0"}
	head := map[string]string{"github.com/foo/baz": "v1.0.0"}
	report := DetectRenames(base, head, "feature")
	if len(report.Entries) != 0 {
		t.Errorf("expected no renames, got %d", len(report.Entries))
	}
}

func TestDetectRenames_Empty(t *testing.T) {
	report := DetectRenames(map[string]string{}, map[string]string{}, "main")
	if len(report.Entries) != 0 {
		t.Errorf("expected empty report, got %d entries", len(report.Entries))
	}
}

func TestFormatRenameReport_NoEntries(t *testing.T) {
	r := RenameReport{Branch: "main", Entries: nil}
	out := FormatRenameReport(r)
	if !strings.Contains(out, "No likely renames") {
		t.Errorf("expected no-renames message, got: %s", out)
	}
}

func TestFormatRenameReport_ContainsModules(t *testing.T) {
	r := RenameReport{
		Branch: "main",
		Entries: []RenameEntry{
			{OldModule: "github.com/old/logger", NewModule: "github.com/new/logger", Version: "v1.2.0", Note: "possible rename"},
		},
	}
	out := FormatRenameReport(r)
	if !strings.Contains(out, "github.com/old/logger") {
		t.Errorf("expected old module in output")
	}
	if !strings.Contains(out, "github.com/new/logger") {
		t.Errorf("expected new module in output")
	}
}

func TestSaveAndLoadRenameReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rename.json")
	r := RenameReport{
		Branch: "main",
		Entries: []RenameEntry{
			{OldModule: "github.com/old/logger", NewModule: "github.com/new/logger", Version: "v1.0.0", Note: "test"},
		},
	}
	if err := SaveRenameReport(path, r); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRenameReport(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Branch != r.Branch {
		t.Errorf("branch mismatch: %s", loaded.Branch)
	}
	if len(loaded.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(loaded.Entries))
	}
}

func TestLoadRenameReport_Missing(t *testing.T) {
	_, err := LoadRenameReport("/nonexistent/path/rename.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRenameReport_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := LoadRenameReport(path)
	if err == nil {
		t.Fatal("expected decode error")
	}
}
