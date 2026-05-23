package gomod

import (
	"testing"
)

func TestFilterDiff_MajorOnly(t *testing.T) {
	input := DiffResult{
		Changed: []DependencyChange{
			{Module: "github.com/foo/bar", OldVersion: "v1.2.3", NewVersion: "v2.0.0"},
			{Module: "github.com/foo/baz", OldVersion: "v1.0.0", NewVersion: "v1.5.0"},
		},
	}
	opts := FilterOptions{MajorOnly: true}
	result := FilterDiff(input, opts)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed dep, got %d", len(result.Changed))
	}
	if result.Changed[0].Module != "github.com/foo/bar" {
		t.Errorf("unexpected module: %s", result.Changed[0].Module)
	}
}

func TestFilterDiff_PrefixFilter(t *testing.T) {
	input := DiffResult{
		Changed: []DependencyChange{
			{Module: "github.com/org/a", OldVersion: "v1.0.0", NewVersion: "v1.1.0"},
			{Module: "golang.org/x/net", OldVersion: "v0.1.0", NewVersion: "v0.2.0"},
		},
	}
	opts := FilterOptions{PrefixFilter: "github.com/org"}
	result := FilterDiff(input, opts)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed dep, got %d", len(result.Changed))
	}
	if result.Changed[0].Module != "github.com/org/a" {
		t.Errorf("unexpected module: %s", result.Changed[0].Module)
	}
}

func TestFilterDiff_IncludeAdded(t *testing.T) {
	input := DiffResult{
		Added:   []Dependency{{Module: "github.com/new/dep", Version: "v1.0.0"}},
		Removed: []Dependency{{Module: "github.com/old/dep", Version: "v1.0.0"}},
	}
	opts := FilterOptions{IncludeAdded: true, IncludeRemoved: false}
	result := FilterDiff(input, opts)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added dep, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Fatalf("expected 0 removed deps, got %d", len(result.Removed))
	}
}

func TestIsMajorChange(t *testing.T) {
	cases := []struct {
		old, new string
		want     bool
	}{
		{"v1.0.0", "v2.0.0", true},
		{"v1.0.0", "v1.9.9", false},
		{"v0.1.0", "v1.0.0", true},
		{"v2.3.1", "v2.4.0", false},
	}
	for _, c := range cases {
		got := isMajorChange(c.old, c.new)
		if got != c.want {
			t.Errorf("isMajorChange(%q, %q) = %v, want %v", c.old, c.new, got, c.want)
		}
	}
}
