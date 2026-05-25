package gomod

import (
	"os"
	"strings"
	"testing"
	"time"
)

var sampleBaseline = []Dependency{
	{Module: "github.com/foo/bar", Version: "v1.0.0"},
	{Module: "github.com/foo/baz", Version: "v2.1.0"},
	{Module: "github.com/stable/lib", Version: "v0.9.0"},
}

var sampleCurrent = []Dependency{
	{Module: "github.com/foo/bar", Version: "v2.0.0"}, // major bump
	{Module: "github.com/foo/baz", Version: "v2.2.0"}, // minor bump
	{Module: "github.com/stable/lib", Version: "v0.9.0"}, // unchanged
}

func TestAnalyzeDrift_EntryCount(t *testing.T) {
	since := time.Now().Add(-30 * 24 * time.Hour)
	report := AnalyzeDrift("main", sampleBaseline, sampleCurrent, since)
	if len(report.Entries) != 2 {
		t.Fatalf("expected 2 drift entries, got %d", len(report.Entries))
	}
}

func TestAnalyzeDrift_MajorIsHigh(t *testing.T) {
	since := time.Now().Add(-10 * 24 * time.Hour)
	report := AnalyzeDrift("main", sampleBaseline, sampleCurrent, since)
	for _, e := range report.Entries {
		if e.Module == "github.com/foo/bar" && e.Severity != "high" {
			t.Errorf("expected high severity for major bump, got %s", e.Severity)
		}
	}
}

func TestAnalyzeDrift_OldMinorIsMedium(t *testing.T) {
	since := time.Now().Add(-100 * 24 * time.Hour)
	report := AnalyzeDrift("main", sampleBaseline, sampleCurrent, since)
	for _, e := range report.Entries {
		if e.Module == "github.com/foo/baz" && e.Severity != "medium" {
			t.Errorf("expected medium severity for old minor bump, got %s", e.Severity)
		}
	}
}

func TestFormatDriftReport_ContainsBranch(t *testing.T) {
	since := time.Now().Add(-5 * 24 * time.Hour)
	report := AnalyzeDrift("feature/x", sampleBaseline, sampleCurrent, since)
	out := FormatDriftReport(report)
	if !strings.Contains(out, "feature/x") {
		t.Errorf("expected branch name in output, got: %s", out)
	}
}

func TestFormatDriftReport_NoDrift(t *testing.T) {
	report := AnalyzeDrift("main", sampleBaseline, sampleBaseline, time.Now())
	out := FormatDriftReport(report)
	if !strings.Contains(out, "no drift") {
		t.Errorf("expected 'no drift' message, got: %s", out)
	}
}

func TestSaveAndLoadDriftReport(t *testing.T) {
	dir := t.TempDir()
	since := time.Now().Add(-20 * 24 * time.Hour)
	report := AnalyzeDrift("main", sampleBaseline, sampleCurrent, since)
	if err := SaveDriftReport(dir, report); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadDriftReport(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != len(report.Entries) {
		t.Errorf("entry count mismatch: want %d got %d", len(report.Entries), len(loaded.Entries))
	}
}

func TestLoadDriftReport_Missing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadDriftReport(dir)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadDriftReport_Invalid(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/drift.json", []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadDriftReport(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
