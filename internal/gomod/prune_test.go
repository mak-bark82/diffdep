package gomod

import (
	"strings"
	"testing"
)

func sampleDepsForPrune() ([]Dependency, []Dependency) {
	base := []Dependency{
		{Module: "github.com/foo/bar", Version: "v1.0.0"},
		{Module: "github.com/foo/baz", Version: "v2.0.0"},
		{Module: "github.com/foo/qux", Version: "v0.5.0"},
	}
	head := []Dependency{
		{Module: "github.com/foo/bar", Version: "v1.0.0"},
		{Module: "github.com/foo/baz", Version: "v3.0.0"},
		// qux removed
	}
	return base, head
}

func TestAnalyzePrune_RemovedDep(t *testing.T) {
	base, head := sampleDepsForPrune()
	diff := DiffDependencies(base, head)
	result := AnalyzePrune(base, head, diff)

	found := false
	for _, s := range result.Suggestions {
		if s.Module == "github.com/foo/qux" {
			found = true
			if !strings.Contains(s.Reason, "removed") {
				t.Errorf("expected 'removed' in reason, got: %s", s.Reason)
			}
		}
	}
	if !found {
		t.Error("expected prune suggestion for removed module github.com/foo/qux")
	}
}

func TestAnalyzePrune_UnchangedDep(t *testing.T) {
	base, head := sampleDepsForPrune()
	diff := DiffDependencies(base, head)
	result := AnalyzePrune(base, head, diff)

	found := false
	for _, s := range result.Suggestions {
		if s.Module == "github.com/foo/bar" {
			found = true
			if !strings.Contains(s.Reason, "unchanged") {
				t.Errorf("expected 'unchanged' in reason, got: %s", s.Reason)
			}
		}
	}
	if !found {
		t.Error("expected prune suggestion for unchanged module github.com/foo/bar")
	}
}

func TestAnalyzePrune_Empty(t *testing.T) {
	result := AnalyzePrune(nil, nil, nil)
	if len(result.Suggestions) != 0 {
		t.Errorf("expected no suggestions, got %d", len(result.Suggestions))
	}
}

func TestFormatPruneResult_NoSuggestions(t *testing.T) {
	out := FormatPruneResult(PruneResult{})
	if !strings.Contains(out, "No prune") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatPruneResult_ContainsModule(t *testing.T) {
	r := PruneResult{
		Suggestions: []PruneSuggestion{
			{Module: "github.com/foo/bar", Version: "v1.0.0", Reason: "removed in head branch"},
		},
	}
	out := FormatPruneResult(r)
	if !strings.Contains(out, "github.com/foo/bar") {
		t.Errorf("expected module name in output, got: %s", out)
	}
	if !strings.Contains(out, "v1.0.0") {
		t.Errorf("expected version in output, got: %s", out)
	}
}
