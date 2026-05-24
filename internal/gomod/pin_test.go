package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForPin() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", OldVersion: "v1.2.0", NewVersion: "v1.3.0", ChangeType: "changed"},
		{Module: "github.com/baz/qux", OldVersion: "v2.0.0", NewVersion: "v2.1.0", ChangeType: "changed"},
		{Module: "github.com/stable/lib", OldVersion: "v0.5.0", NewVersion: "v0.5.0", ChangeType: "unchanged"},
	}
}

func TestCheckViolations_NoViolations(t *testing.T) {
	pl := NewPinList()
	pl.Add("github.com/foo/bar", "v1.3.0", "")

	violations := pl.CheckViolations(sampleDiffForPin())
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckViolations_DetectsViolation(t *testing.T) {
	pl := NewPinList()
	pl.Add("github.com/foo/bar", "v1.2.0", "security fix")

	violations := pl.CheckViolations(sampleDiffForPin())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	v := violations[0]
	if v.Module != "github.com/foo/bar" {
		t.Errorf("unexpected module: %s", v.Module)
	}
	if v.Pinned != "v1.2.0" || v.Actual != "v1.3.0" {
		t.Errorf("unexpected versions pinned=%s actual=%s", v.Pinned, v.Actual)
	}
	if v.Reason != "security fix" {
		t.Errorf("unexpected reason: %s", v.Reason)
	}
}

func TestCheckViolations_MultipleViolations(t *testing.T) {
	pl := NewPinList()
	pl.Add("github.com/foo/bar", "v1.0.0", "")
	pl.Add("github.com/baz/qux", "v2.0.0", "locked")

	violations := pl.CheckViolations(sampleDiffForPin())
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestFormatPinViolations_Empty(t *testing.T) {
	out := FormatPinViolations(nil)
	if out != "No pin violations found." {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatPinViolations_ContainsModule(t *testing.T) {
	violations := []PinViolation{
		{Module: "github.com/foo/bar", Pinned: "v1.0.0", Actual: "v1.3.0", Reason: "audit"},
	}
	out := FormatPinViolations(violations)
	if !strings.Contains(out, "github.com/foo/bar") {
		t.Errorf("expected module in output, got: %s", out)
	}
	if !strings.Contains(out, "audit") {
		t.Errorf("expected reason in output, got: %s", out)
	}
	if !strings.Contains(out, "v1.0.0") || !strings.Contains(out, "v1.3.0") {
		t.Errorf("expected versions in output, got: %s", out)
	}
}
