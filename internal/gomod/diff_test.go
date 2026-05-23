package gomod

import (
	"testing"
)

func TestDepsToMap(t *testing.T) {
	deps := []Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
		{Path: "github.com/baz/qux", Version: "v2.3.1"},
	}

	m := DepsToMap(deps)

	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
	if m["github.com/foo/bar"] != "v1.0.0" {
		t.Errorf("unexpected version for foo/bar: %s", m["github.com/foo/bar"])
	}
	if m["github.com/baz/qux"] != "v2.3.1" {
		t.Errorf("unexpected version for baz/qux: %s", m["github.com/baz/qux"])
	}
}

func TestDiffDependencies_Added(t *testing.T) {
	base := []Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
	}
	head := []Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
		{Path: "github.com/new/pkg", Version: "v0.1.0"},
	}

	result := DiffDependencies(base, head)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(result.Added))
	}
	if result.Added[0].Path != "github.com/new/pkg" {
		t.Errorf("unexpected added path: %s", result.Added[0].Path)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
	if len(result.Changed) != 0 {
		t.Errorf("expected 0 changed, got %d", len(result.Changed))
	}
}

func TestDiffDependencies_Removed(t *testing.T) {
	base := []Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
		{Path: "github.com/old/pkg", Version: "v3.0.0"},
	}
	head := []Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
	}

	result := DiffDependencies(base, head)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(result.Removed))
	}
	if result.Removed[0].Path != "github.com/old/pkg" {
		t.Errorf("unexpected removed path: %s", result.Removed[0].Path)
	}
}

func TestDiffDependencies_Changed(t *testing.T) {
	base := []Dependency{
		{Path: "github.com/foo/bar", Version: "v1.0.0"},
	}
	head := []Dependency{
		{Path: "github.com/foo/bar", Version: "v2.0.0"},
	}

	result := DiffDependencies(base, head)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(result.Changed))
	}
	change := result.Changed[0]
	if change.Path != "github.com/foo/bar" {
		t.Errorf("unexpected changed path: %s", change.Path)
	}
	if change.OldVersion != "v1.0.0" {
		t.Errorf("unexpected old version: %s", change.OldVersion)
	}
	if change.NewVersion != "v2.0.0" {
		t.Errorf("unexpected new version: %s", change.NewVersion)
	}
}

func TestDiffDependencies_Empty(t *testing.T) {
	result := DiffDependencies(nil, nil)

	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Error("expected empty diff for nil inputs")
	}
}
