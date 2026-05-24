package gomod

import (
	"strings"
	"testing"
)

func sampleDepsForRetract() []Dependency {
	return []Dependency{
		{Module: "github.com/pkg/errors", Version: "v0.8.0"},
		{Module: "github.com/pkg/errors", Version: "v0.9.1"},
		{Module: "github.com/gogo/protobuf", Version: "v1.3.1"},
		{Module: "golang.org/x/net", Version: "v0.5.0"},
	}
}

func TestCheckRetractions_DetectsKnown(t *testing.T) {
	deps := sampleDepsForRetract()
	report := CheckRetractions("main", deps)
	if len(report.Retracted) != 2 {
		t.Fatalf("expected 2 retracted, got %d", len(report.Retracted))
	}
}

func TestCheckRetractions_SkipsSafeVersion(t *testing.T) {
	deps := []Dependency{
		{Module: "github.com/pkg/errors", Version: "v0.9.1"},
	}
	report := CheckRetractions("main", deps)
	if len(report.Retracted) != 0 {
		t.Errorf("expected no retractions, got %d", len(report.Retracted))
	}
}

func TestCheckRetractions_Empty(t *testing.T) {
	report := CheckRetractions("feature", []Dependency{})
	if len(report.Retracted) != 0 {
		t.Errorf("expected no retractions for empty deps")
	}
	if report.Branch != "feature" {
		t.Errorf("expected branch 'feature', got %s", report.Branch)
	}
}

func TestFormatRetractReport_NoRetractions(t *testing.T) {
	report := RetractReport{Branch: "main", Retracted: nil}
	out := FormatRetractReport(report)
	if !strings.Contains(out, "No retracted modules") {
		t.Errorf("expected no-retraction message, got: %s", out)
	}
}

func TestFormatRetractReport_ContainsModule(t *testing.T) {
	report := RetractReport{
		Branch: "main",
		Retracted: []RetractedModule{
			{Module: "github.com/pkg/errors", Version: "v0.8.0", Reason: "use v0.9.1+"},
		},
	}
	out := FormatRetractReport(report)
	if !strings.Contains(out, "github.com/pkg/errors") {
		t.Errorf("expected module name in output, got: %s", out)
	}
	if !strings.Contains(out, "use v0.9.1+") {
		t.Errorf("expected reason in output, got: %s", out)
	}
	if !strings.Contains(out, "RETRACTED") {
		t.Errorf("expected RETRACTED label in output")
	}
}
