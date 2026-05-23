package gomod

import (
	"testing"
)

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 || s.Added != 0 || s.Removed != 0 || s.Changed != 0 {
		t.Errorf("expected zero summary, got %+v", s)
	}
	if s.HasChanges() {
		t.Error("expected HasChanges() == false for empty diff")
	}
}

func TestSummarize_Mixed(t *testing.T) {
	entries := []DiffEntry{
		{Module: "github.com/a/b", OldVersion: "", NewVersion: "v1.0.0", ChangeType: Added},
		{Module: "github.com/c/d", OldVersion: "v2.0.0", NewVersion: "", ChangeType: Removed},
		{Module: "github.com/e/f", OldVersion: "v1.0.0", NewVersion: "v2.0.0", ChangeType: Changed},
		{Module: "github.com/g/h", OldVersion: "v1.1.0", NewVersion: "v1.2.0", ChangeType: Changed},
	}
	s := Summarize(entries)
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Changed != 2 {
		t.Errorf("expected Changed=2, got %d", s.Changed)
	}
	if !s.HasChanges() {
		t.Error("expected HasChanges() == true")
	}
}

func TestSummary_String(t *testing.T) {
	s := Summary{Total: 3, Added: 1, Removed: 1, Changed: 1}
	got := s.String()
	want := "total=3 added=1 removed=1 changed=1"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
