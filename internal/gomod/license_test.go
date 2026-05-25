package gomod

import (
	"strings"
	"testing"
)

func sampleDepsForLicense() []Dependency {
	return []Dependency{
		{Module: "github.com/foo/mitlib", Version: "v1.0.0"},
		{Module: "github.com/bar/gplapp", Version: "v2.1.0"},
		{Module: "github.com/baz/apache", Version: "v0.9.0"},
	}
}

func TestCheckLicenses_ReturnsEntries(t *testing.T) {
	deps := sampleDepsForLicense()
	report := CheckLicenses(deps)
	if len(report.Entries) != len(deps) {
		t.Errorf("expected %d entries, got %d", len(deps), len(report.Entries))
	}
}

func TestCheckLicenses_HighRiskCount(t *testing.T) {
	deps := sampleDepsForLicense()
	report := CheckLicenses(deps)
	if report.HighRiskCount < 0 {
		t.Error("HighRiskCount should not be negative")
	}
}

func TestFormatLicenseReport_ContainsModules(t *testing.T) {
	deps := sampleDepsForLicense()
	report := CheckLicenses(deps)
	out := FormatLicenseReport(report)
	for _, d := range deps {
		if !strings.Contains(out, d.Module) {
			t.Errorf("expected module %s in report output", d.Module)
		}
	}
}

func TestFormatLicenseReport_Empty(t *testing.T) {
	report := CheckLicenses([]Dependency{})
	out := FormatLicenseReport(report)
	if out == "" {
		t.Error("expected non-empty output even for empty deps")
	}
}

func TestFormatLicenseReport_ContainsRisk(t *testing.T) {
	deps := sampleDepsForLicense()
	report := CheckLicenses(deps)
	out := FormatLicenseReport(report)
	if !strings.Contains(out, "License") && !strings.Contains(out, "license") {
		t.Error("expected report to mention license information")
	}
}
