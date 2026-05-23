package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadAlerts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")

	alerts := []Alert{
		{Module: "github.com/foo/bar", Level: AlertCritical, Message: "major version change"},
		{Module: "github.com/baz/qux", Level: AlertWarning, Message: "dependency removed"},
	}

	if err := SaveAlerts(path, alerts); err != nil {
		t.Fatalf("SaveAlerts failed: %v", err)
	}

	loaded, err := LoadAlerts(path)
	if err != nil {
		t.Fatalf("LoadAlerts failed: %v", err)
	}

	if len(loaded) != len(alerts) {
		t.Fatalf("expected %d alerts, got %d", len(alerts), len(loaded))
	}
	for i, a := range loaded {
		if a.Module != alerts[i].Module {
			t.Errorf("alert[%d] module mismatch: got %s", i, a.Module)
		}
		if a.Level != alerts[i].Level {
			t.Errorf("alert[%d] level mismatch: got %s", i, a.Level)
		}
	}
}

func TestLoadAlerts_Missing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")

	alerts, err := LoadAlerts(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(alerts) != 0 {
		t.Errorf("expected empty slice, got %d alerts", len(alerts))
	}
}

func TestLoadAlerts_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")

	if err := os.WriteFile(path, []byte("not-json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadAlerts(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
