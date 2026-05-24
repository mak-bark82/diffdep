package gomod

import (
	"os"
	"testing"
)

func sampleDepsForLock() []Dependency {
	return []Dependency{
		{Module: "github.com/foo/bar", Version: "v1.2.3"},
		{Module: "github.com/baz/qux", Version: "v2.0.0"},
	}
}

func TestNewLockFile_SortsEntries(t *testing.T) {
	deps := sampleDepsForLock()
	lf := NewLockFile(deps)
	if len(lf.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lf.Entries))
	}
	if lf.Entries[0].Module != "github.com/baz/qux" {
		t.Errorf("expected sorted first entry, got %s", lf.Entries[0].Module)
	}
}

func TestCheckLock_NoViolations(t *testing.T) {
	lf := NewLockFile(sampleDepsForLock())
	diff := []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.2.3", NewVersion: "v1.2.3", ChangeType: "unchanged"},
	}
	violations := CheckLock(lf, diff)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckLock_DetectsViolation(t *testing.T) {
	lf := NewLockFile(sampleDepsForLock())
	diff := []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.2.3", NewVersion: "v1.3.0", ChangeType: "changed"},
	}
	violations := CheckLock(lf, diff)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0] == "" {
		t.Error("expected non-empty violation message")
	}
}

func TestCheckLock_NilLockFile(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.2.3", NewVersion: "v1.3.0", ChangeType: "changed"},
	}
	violations := CheckLock(nil, diff)
	if violations != nil {
		t.Errorf("expected nil violations for nil lock file")
	}
}

func TestFormatLockViolations_Empty(t *testing.T) {
	out := FormatLockViolations(nil)
	if out != "No lock violations found.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatLockViolations_WithViolations(t *testing.T) {
	violations := []string{"github.com/foo/bar: locked to v1.2.3, got v1.3.0"}
	out := FormatLockViolations(violations)
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestSaveAndLoadLockFile(t *testing.T) {
	dir := t.TempDir()
	lf := NewLockFile(sampleDepsForLock())
	if err := SaveLockFile(dir, lf); err != nil {
		t.Fatalf("SaveLockFile: %v", err)
	}
	loaded, err := LoadLockFile(dir)
	if err != nil {
		t.Fatalf("LoadLockFile: %v", err)
	}
	if len(loaded.Entries) != len(lf.Entries) {
		t.Errorf("entry count mismatch: want %d, got %d", len(lf.Entries), len(loaded.Entries))
	}
}

func TestLoadLockFile_Missing(t *testing.T) {
	dir := t.TempDir()
	lf, err := LoadLockFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf != nil {
		t.Error("expected nil for missing lock file")
	}
}

func TestLoadLockFile_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/diffdep.lock.json"
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadLockFile(dir)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
