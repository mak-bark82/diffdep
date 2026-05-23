package gomod_test

import (
	"path/filepath"
	"testing"

	"github.com/user/diffdep/internal/gomod"
)

func TestCompareAgainstBaseline_NoChanges(t *testing.T) {
	deps := []gomod.Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
	}
	tmp := t.TempDir()
	path := filepath.Join(tmp, "baseline.json")
	if err := gomod.SaveBaseline(path, "main", deps); err != nil {
		t.Fatal(err)
	}

	cmp, err := gomod.CompareAgainstBaseline(path, deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmp.Diff) != 0 {
		t.Errorf("expected no diff, got %d entries", len(cmp.Diff))
	}
	if cmp.HasBreakingChanges() {
		t.Error("expected no breaking changes")
	}
}

func TestCompareAgainstBaseline_Breaking(t *testing.T) {
	old := []gomod.Dependency{{Path: "github.com/foo/bar", Version: "v1.3.0"}}
	current := []gomod.Dependency{{Path: "github.com/foo/bar", Version: "v2.0.0"}}

	tmp := t.TempDir()
	path := filepath.Join(tmp, "baseline.json")
	if err := gomod.SaveBaseline(path, "main", old); err != nil {
		t.Fatal(err)
	}

	cmp, err := gomod.CompareAgainstBaseline(path, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cmp.HasBreakingChanges() {
		t.Error("expected breaking changes")
	}
}

func TestCompareAgainstBaseline_MissingFile(t *testing.T) {
	_, err := gomod.CompareAgainstBaseline("/no/such/file.json", nil)
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}
