package gomod

import (
	"strings"
	"testing"
)

var baseOverlap = []Dependency{
	{Module: "github.com/foo/bar", Version: "v1.2.0"},
	{Module: "github.com/foo/baz", Version: "v2.0.0"},
	{Module: "github.com/only/base", Version: "v0.1.0"},
}

var headOverlap = []Dependency{
	{Module: "github.com/foo/bar", Version: "v1.2.0"},
	{Module: "github.com/foo/baz", Version: "v3.0.0"},
	{Module: "github.com/only/head", Version: "v0.2.0"},
}

func TestAnalyzeOverlap_SharedCount(t *testing.T) {
	r := AnalyzeOverlap("main", baseOverlap, headOverlap)
	if r.Shared != 1 {
		t.Errorf("expected 1 shared, got %d", r.Shared)
	}
}

func TestAnalyzeOverlap_DivergeCount(t *testing.T) {
	r := AnalyzeOverlap("main", baseOverlap, headOverlap)
	if r.Diverge != 1 {
		t.Errorf("expected 1 diverged, got %d", r.Diverge)
	}
}

func TestAnalyzeOverlap_ExcludesUnique(t *testing.T) {
	r := AnalyzeOverlap("main", baseOverlap, headOverlap)
	for _, e := range r.Entries {
		if e.Module == "github.com/only/base" || e.Module == "github.com/only/head" {
			t.Errorf("unique module %s should not appear in overlap", e.Module)
		}
	}
}

func TestAnalyzeOverlap_Empty(t *testing.T) {
	r := AnalyzeOverlap("main", nil, nil)
	if r.Shared != 0 || r.Diverge != 0 || len(r.Entries) != 0 {
		t.Error("expected empty report for nil inputs")
	}
}

func TestAnalyzeOverlap_BranchSet(t *testing.T) {
	r := AnalyzeOverlap("feature-x", baseOverlap, headOverlap)
	if r.Branch != "feature-x" {
		t.Errorf("expected branch feature-x, got %s", r.Branch)
	}
}

func TestFormatOverlapReport_ContainsBranch(t *testing.T) {
	r := AnalyzeOverlap("main", baseOverlap, headOverlap)
	out := FormatOverlapReport(r)
	if !strings.Contains(out, "main") {
		t.Error("expected branch name in output")
	}
}

func TestFormatOverlapReport_NoShared(t *testing.T) {
	base := []Dependency{{Module: "github.com/a/a", Version: "v1.0.0"}}
	head := []Dependency{{Module: "github.com/b/b", Version: "v1.0.0"}}
	r := AnalyzeOverlap("main", base, head)
	out := FormatOverlapReport(r)
	if !strings.Contains(out, "No shared modules") {
		t.Error("expected no-shared message")
	}
}

func TestFormatOverlapReport_DivergeMark(t *testing.T) {
	r := AnalyzeOverlap("main", baseOverlap, headOverlap)
	out := FormatOverlapReport(r)
	if !strings.Contains(out, "[~]") {
		t.Error("expected diverge marker [~] in output")
	}
}
