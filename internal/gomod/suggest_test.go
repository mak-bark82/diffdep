package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForSuggest() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "", NewVersion: "v1.0.0", Change: ChangeAdded},
		{Module: "github.com/foo/baz", OldVersion: "v1.3.0", NewVersion: "", Change: ChangeRemoved},
		{Module: "github.com/foo/qux", OldVersion: "v1.2.0", NewVersion: "v1.3.0", Change: ChangeUpdated},
		{Module: "github.com/foo/major", OldVersion: "v1.0.0", NewVersion: "v2.0.0", Change: ChangeUpdated},
		{Module: "github.com/foo/down", OldVersion: "v2.0.0", NewVersion: "v1.9.0", Change: ChangeUpdated},
	}
}

func TestGenerateSuggestions_Empty(t *testing.T) {
	suggestions := GenerateSuggestions([]DiffEntry{})
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(suggestions))
	}
}

func TestGenerateSuggestions_Added(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/bar", NewVersion: "v1.0.0", Change: ChangeAdded},
	}
	suggestions := GenerateSuggestions(diff)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Priority != "low" {
		t.Errorf("expected low priority for added dep, got %s", suggestions[0].Priority)
	}
}

func TestGenerateSuggestions_Removed(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/baz", OldVersion: "v1.3.0", Change: ChangeRemoved},
	}
	suggestions := GenerateSuggestions(diff)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Priority != "medium" {
		t.Errorf("expected medium priority for removed dep, got %s", suggestions[0].Priority)
	}
}

func TestGenerateSuggestions_MajorBump(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/major", OldVersion: "v1.0.0", NewVersion: "v2.0.0", Change: ChangeUpdated},
	}
	suggestions := GenerateSuggestions(diff)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Priority != "high" {
		t.Errorf("expected high priority for major bump, got %s", suggestions[0].Priority)
	}
}

func TestGenerateSuggestions_Downgrade(t *testing.T) {
	diff := []DiffEntry{
		{Module: "github.com/foo/down", OldVersion: "v2.0.0", NewVersion: "v1.9.0", Change: ChangeUpdated},
	}
	suggestions := GenerateSuggestions(diff)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Priority != "medium" {
		t.Errorf("expected medium priority for downgrade, got %s", suggestions[0].Priority)
	}
}

func TestFormatSuggestions_NoSuggestions(t *testing.T) {
	out := FormatSuggestions([]Suggestion{})
	if out != "No suggestions.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatSuggestions_ContainsModule(t *testing.T) {
	suggestions := GenerateSuggestions(sampleDiffForSuggest())
	out := FormatSuggestions(suggestions)
	if !strings.Contains(out, "github.com/foo/major") {
		t.Errorf("expected output to contain major module, got:\n%s", out)
	}
	if !strings.Contains(out, "Suggestions:") {
		t.Errorf("expected header in output")
	}
}
