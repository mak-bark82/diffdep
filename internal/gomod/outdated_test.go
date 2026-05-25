package gomod

import (
	"errors"
	"strings"
	"testing"
)

func sampleDepsForOutdated() []Dependency {
	return []Dependency{
		{Module: "github.com/foo/bar", Version: "v1.2.0"},
		{Module: "github.com/baz/qux", Version: "v2.0.0"},
		{Module: "github.com/up/todate", Version: "v3.1.0"},
	}
}

func stubResolver(latest map[string]string) LatestResolver {
	return func(module string) (string, error) {
		if v, ok := latest[module]; ok {
			return v, nil
		}
		return "", errors.New("unknown module")
	}
}

func TestCheckOutdated_DetectsStale(t *testing.T) {
	deps := sampleDepsForOutdated()
	resolver := stubResolver(map[string]string{
		"github.com/foo/bar":  "v1.3.0",
		"github.com/baz/qux":  "v2.0.0",
		"github.com/up/todate": "v3.1.0",
	})
	report, err := CheckOutdated("main", deps, resolver)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Branch != "main" {
		t.Errorf("expected branch main, got %s", report.Branch)
	}
	var stale []OutdatedEntry
	for _, e := range report.Entries {
		if e.IsStale {
			stale = append(stale, e)
		}
	}
	if len(stale) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(stale))
	}
	if stale[0].Module != "github.com/foo/bar" {
		t.Errorf("unexpected stale module: %s", stale[0].Module)
	}
}

func TestCheckOutdated_SkipsUnresolvable(t *testing.T) {
	deps := sampleDepsForOutdated()
	// resolver returns error for all
	resolver := stubResolver(map[string]string{})
	report, err := CheckOutdated("feature", deps, resolver)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(report.Entries))
	}
}

func TestCheckOutdated_AllCurrent(t *testing.T) {
	deps := sampleDepsForOutdated()
	resolver := stubResolver(map[string]string{
		"github.com/foo/bar":  "v1.2.0",
		"github.com/baz/qux":  "v2.0.0",
		"github.com/up/todate": "v3.1.0",
	})
	report, _ := CheckOutdated("main", deps, resolver)
	for _, e := range report.Entries {
		if e.IsStale {
			t.Errorf("expected no stale entries, got %s", e.Module)
		}
	}
}

func TestFormatOutdatedReport_ContainsBranch(t *testing.T) {
	report := &OutdatedReport{
		Branch: "develop",
		Entries: []OutdatedEntry{
			{Module: "github.com/foo/bar", Current: "v1.0.0", Latest: "v1.1.0", IsStale: true},
		},
	}
	out := FormatOutdatedReport(report)
	if !strings.Contains(out, "develop") {
		t.Errorf("expected branch name in output, got: %s", out)
	}
	if !strings.Contains(out, "github.com/foo/bar") {
		t.Errorf("expected module in output, got: %s", out)
	}
}

func TestFormatOutdatedReport_Empty(t *testing.T) {
	report := &OutdatedReport{Branch: "main", Entries: []OutdatedEntry{}}
	out := FormatOutdatedReport(report)
	if !strings.Contains(out, "up to date") {
		t.Errorf("expected up-to-date message, got: %s", out)
	}
}
