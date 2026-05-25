package gomod

import (
	"os"
	"testing"
)

func sampleLicenseReport() *LicenseReport {
	return &LicenseReport{
		Entries: []LicenseEntry{
			{Module: "github.com/foo/bar", Version: "v1.2.3", License: "MIT", Risk: "low"},
			{Module: "github.com/baz/qux", Version: "v2.0.0", License: "GPL-3.0", Risk: "high"},
		},
		HighRiskCount: 1,
	}
}

func TestSaveAndLoadLicenseReport(t *testing.T) {
	dir := t.TempDir()
	report := sampleLicenseReport()

	if err := SaveLicenseReport(dir, report); err != nil {
		t.Fatalf("SaveLicenseReport: %v", err)
	}

	loaded, err := LoadLicenseReport(dir)
	if err != nil {
		t.Fatalf("LoadLicenseReport: %v", err)
	}

	if len(loaded.Entries) != len(report.Entries) {
		t.Errorf("expected %d entries, got %d", len(report.Entries), len(loaded.Entries))
	}
	if loaded.HighRiskCount != report.HighRiskCount {
		t.Errorf("expected HighRiskCount %d, got %d", report.HighRiskCount, loaded.HighRiskCount)
	}
	if loaded.Entries[0].Module != "github.com/foo/bar" {
		t.Errorf("unexpected first module: %s", loaded.Entries[0].Module)
	}
}

func TestLoadLicenseReport_Missing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadLicenseReport(dir)
	if err == nil {
		t.Fatal("expected error for missing report, got nil")
	}
}

func TestLoadLicenseReport_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/license_report.json"
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadLicenseReport(dir)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
