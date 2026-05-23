package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForAlert() DiffResult {
	return DiffResult{
		Added: []Dependency{
			{Module: "github.com/new/pkg", Version: "v1.0.0"},
		},
		Removed: []Dependency{
			{Module: "github.com/old/pkg", Version: "v2.1.0"},
		},
		Changed: []DependencyChange{
			{Module: "github.com/foo/bar", OldVersion: "v1.2.0", NewVersion: "v2.0.0"},
			{Module: "github.com/baz/qux", OldVersion: "v1.0.0", NewVersion: "v1.3.0"},
		},
	}
}

func TestGenerateAlerts_DefaultConfig(t *testing.T) {
	diff := sampleDiffForAlert()
	cfg := DefaultAlertConfig()
	alerts := GenerateAlerts(diff, cfg)

	// default: OnRemoved + OnMajor only
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestGenerateAlerts_AllEnabled(t *testing.T) {
	diff := sampleDiffForAlert()
	cfg := AlertConfig{OnAdded: true, OnRemoved: true, OnMajor: true, OnMinor: true}
	alerts := GenerateAlerts(diff, cfg)

	if len(alerts) != 4 {
		t.Fatalf("expected 4 alerts, got %d", len(alerts))
	}
}

func TestGenerateAlerts_Empty(t *testing.T) {
	alerts := GenerateAlerts(DiffResult{}, DefaultAlertConfig())
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestGenerateAlerts_CriticalLevel(t *testing.T) {
	diff := DiffResult{
		Changed: []DependencyChange{
			{Module: "github.com/foo/bar", OldVersion: "v1.0.0", NewVersion: "v2.0.0"},
		},
	}
	cfg := AlertConfig{OnMajor: true}
	alerts := GenerateAlerts(diff, cfg)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != AlertCritical {
		t.Errorf("expected CRITICAL, got %s", alerts[0].Level)
	}
}

func TestFormatAlerts_NoAlerts(t *testing.T) {
	out := FormatAlerts(nil)
	if out != "no alerts" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatAlerts_ContainsLevel(t *testing.T) {
	alerts := []Alert{
		{Module: "mod", Level: AlertWarning, Message: "dependency removed: mod@v1.0.0"},
	}
	out := FormatAlerts(alerts)
	if !strings.Contains(out, "[WARNING]") {
		t.Errorf("expected [WARNING] in output, got: %s", out)
	}
	if !strings.Contains(out, "mod@v1.0.0") {
		t.Errorf("expected module info in output, got: %s", out)
	}
}
