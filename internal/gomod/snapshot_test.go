package gomod_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/diffdep/internal/gomod"
)

func sampleDeps() []gomod.Dependency {
	return []gomod.Dependency{
		{Module: "github.com/foo/bar", Version: "v1.2.3"},
		{Module: "github.com/baz/qux", Version: "v2.0.0"},
	}
}

func TestNewSnapshot_SortsDeps(t *testing.T) {
	deps := sampleDeps()
	s := gomod.NewSnapshot("main", deps)
	if s.Deps[0].Module != "github.com/baz/qux" {
		t.Errorf("expected sorted first dep, got %s", s.Deps[0].Module)
	}
	if s.Branch != "main" {
		t.Errorf("expected branch main, got %s", s.Branch)
	}
}

func TestDiffSnapshot_DetectsChange(t *testing.T) {
	base := gomod.NewSnapshot("main", sampleDeps())
	headDeps := []gomod.Dependency{
		{Module: "github.com/foo/bar", Version: "v1.3.0"},
		{Module: "github.com/baz/qux", Version: "v2.0.0"},
	}
	head := gomod.NewSnapshot("feature", headDeps)
	diff := gomod.DiffSnapshot(base, head)
	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed dep, got %d", len(diff.Changed))
	}
}

func TestSnapshotSummary_ContainsBranch(t *testing.T) {
	s := gomod.NewSnapshot("release/v2", sampleDeps())
	summary := gomod.SnapshotSummary(s)
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	if !containsStr(summary, "release/v2") {
		t.Errorf("summary missing branch name: %s", summary)
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	s := gomod.NewSnapshot("main", sampleDeps())
	if err := gomod.SaveSnapshot(dir, s); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	loaded, err := gomod.LoadSnapshot(dir, "main")
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if len(loaded.Deps) != len(s.Deps) {
		t.Errorf("dep count mismatch: got %d want %d", len(loaded.Deps), len(s.Deps))
	}
}

func TestLoadSnapshot_Missing(t *testing.T) {
	dir := t.TempDir()
	_, err := gomod.LoadSnapshot(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestSaveSnapshot_SlashBranch(t *testing.T) {
	dir := t.TempDir()
	s := gomod.NewSnapshot("feature/my-feature", sampleDeps())
	if err := gomod.SaveSnapshot(dir, s); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	name := entries[0].Name()
	expected := filepath.Base("feature_my-feature.snapshot.json")
	if name != expected {
		t.Errorf("expected filename %s, got %s", expected, name)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
