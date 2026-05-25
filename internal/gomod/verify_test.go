package gomod

import (
	"strings"
	"testing"
)

func sampleDepsForVerify() []Dependency {
	return []Dependency{
		{Module: "github.com/foo/bar", Version: "v1.2.3"},
		{Module: "github.com/baz/qux", Version: "v2.0.0"},
		{Module: "github.com/missing/pkg", Version: "v0.1.0"},
	}
}

func TestVerifyDeps_AllOk(t *testing.T) {
	deps := sampleDepsForVerify()
	checksums := map[string]string{
		"github.com/foo/bar@v1.2.3": "h1:abc123",
		"github.com/baz/qux@v2.0.0": "h1:def456",
		"github.com/missing/pkg@v0.1.0": "h1:ghi789",
	}
	report := VerifyDeps("main", deps, checksums)
	if report.Failed != 0 {
		t.Errorf("expected 0 failures, got %d", report.Failed)
	}
	for _, r := range report.Results {
		if r.Status != "ok" {
			t.Errorf("expected ok for %s, got %s", r.Module, r.Status)
		}
	}
}

func TestVerifyDeps_MissingChecksum(t *testing.T) {
	deps := sampleDepsForVerify()
	checksums := map[string]string{
		"github.com/foo/bar@v1.2.3": "h1:abc123",
		// baz/qux and missing/pkg are absent
	}
	report := VerifyDeps("feature", deps, checksums)
	if report.Failed != 2 {
		t.Errorf("expected 2 failures, got %d", report.Failed)
	}
}

func TestVerifyDeps_BadFormat(t *testing.T) {
	deps := []Dependency{
		{Module: "github.com/weird/hash", Version: "v1.0.0"},
	}
	checksums := map[string]string{
		"github.com/weird/hash@v1.0.0": "sha256:notvalid",
	}
	report := VerifyDeps("main", deps, checksums)
	if report.Failed != 1 {
		t.Errorf("expected 1 failure, got %d", report.Failed)
	}
	if report.Results[0].Status != "checksum_mismatch" {
		t.Errorf("expected checksum_mismatch, got %s", report.Results[0].Status)
	}
}

func TestFormatVerifyReport_ContainsBranch(t *testing.T) {
	report := VerifyReport{Branch: "develop", Failed: 0}
	out := FormatVerifyReport(report)
	if !strings.Contains(out, "develop") {
		t.Error("expected branch name in output")
	}
}

func TestFormatVerifyReport_ContainsModules(t *testing.T) {
	deps := sampleDepsForVerify()
	checksums := map[string]string{
		"github.com/foo/bar@v1.2.3": "h1:abc",
		"github.com/baz/qux@v2.0.0": "h1:def",
		"github.com/missing/pkg@v0.1.0": "h1:ghi",
	}
	report := VerifyDeps("main", deps, checksums)
	out := FormatVerifyReport(report)
	if !strings.Contains(out, "github.com/foo/bar") {
		t.Error("expected module name in output")
	}
}
