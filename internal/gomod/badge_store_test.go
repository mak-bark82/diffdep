package gomod

import (
	"os"
	"testing"
)

func TestSaveAndLoadBadge(t *testing.T) {
	dir := t.TempDir()

	orig := BadgeData{
		Label:   "dep-risk",
		Message: "low (5)",
		Color:   "4c1",
		Style:   BadgeStyleFlat,
	}

	if err := SaveBadge(dir, orig); err != nil {
		t.Fatalf("SaveBadge: %v", err)
	}

	loaded, err := LoadBadge(dir)
	if err != nil {
		t.Fatalf("LoadBadge: %v", err)
	}

	if loaded.Label != orig.Label {
		t.Errorf("label mismatch: got %q want %q", loaded.Label, orig.Label)
	}
	if loaded.Message != orig.Message {
		t.Errorf("message mismatch: got %q want %q", loaded.Message, orig.Message)
	}
	if loaded.Color != orig.Color {
		t.Errorf("color mismatch: got %q want %q", loaded.Color, orig.Color)
	}
	if loaded.Style != orig.Style {
		t.Errorf("style mismatch: got %q want %q", loaded.Style, orig.Style)
	}
}

func TestLoadBadge_Missing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadBadge(dir)
	if err == nil {
		t.Fatal("expected error for missing badge file")
	}
}

func TestLoadBadge_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/badge.json"
	if err := os.WriteFile(path, []byte("not-json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadBadge(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
