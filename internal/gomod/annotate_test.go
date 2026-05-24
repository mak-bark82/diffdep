package gomod

import (
	"strings"
	"testing"
)

func sampleDiffForAnnotate() DiffResult {
	return DiffResult{
		Added: []Dependency{
			{Module: "github.com/new/lib", Version: "v1.0.0"},
		},
		Removed: []Dependency{
			{Module: "github.com/old/lib", Version: "v2.3.1"},
		},
		Changed: []DiffEntry{
			{Module: "github.com/foo/bar", OldVersion: "v1.2.0", NewVersion: "v2.0.0"},
			{Module: "github.com/foo/baz", OldVersion: "v1.5.0", NewVersion: "v1.6.0"},
			{Module: "github.com/foo/qux", OldVersion: "v1.4.0", NewVersion: "v1.3.0"},
		},
	}
}

func TestAnnotateDiff_NewEntry(t *testing.T) {
	diff := sampleDiffForAnnotate()
	annotations := AnnotateDiff(diff)

	var found bool
	for _, a := range annotations {
		if a.Module == "github.com/new/lib" && a.Kind == AnnotationNew {
			found = true
		}
	}
	if !found {
		t.Error("expected new annotation for github.com/new/lib")
	}
}

func TestAnnotateDiff_RemovedEntry(t *testing.T) {
	diff := sampleDiffForAnnotate()
	annotations := AnnotateDiff(diff)

	var found bool
	for _, a := range annotations {
		if a.Module == "github.com/old/lib" && a.Kind == AnnotationRemoved {
			found = true
		}
	}
	if !found {
		t.Error("expected removed annotation for github.com/old/lib")
	}
}

func TestAnnotateDiff_BreakingChange(t *testing.T) {
	diff := sampleDiffForAnnotate()
	annotations := AnnotateDiff(diff)

	var found bool
	for _, a := range annotations {
		if a.Module == "github.com/foo/bar" && a.Kind == AnnotationBreaking {
			found = true
		}
	}
	if !found {
		t.Error("expected breaking annotation for github.com/foo/bar")
	}
}

func TestAnnotateDiff_Downgrade(t *testing.T) {
	diff := sampleDiffForAnnotate()
	annotations := AnnotateDiff(diff)

	var found bool
	for _, a := range annotations {
		if a.Module == "github.com/foo/qux" && a.Kind == AnnotationDowngrade {
			found = true
		}
	}
	if !found {
		t.Error("expected downgrade annotation for github.com/foo/qux")
	}
}

func TestAnnotateDiff_Upgrade(t *testing.T) {
	diff := sampleDiffForAnnotate()
	annotations := AnnotateDiff(diff)

	var found bool
	for _, a := range annotations {
		if a.Module == "github.com/foo/baz" && a.Kind == AnnotationUpgrade {
			found = true
		}
	}
	if !found {
		t.Error("expected upgrade annotation for github.com/foo/baz")
	}
}

func TestFormatAnnotations_Empty(t *testing.T) {
	out := FormatAnnotations(nil)
	if !strings.Contains(out, "no annotations") {
		t.Errorf("expected 'no annotations', got %q", out)
	}
}

func TestFormatAnnotations_ContainsModule(t *testing.T) {
	diff := sampleDiffForAnnotate()
	annotations := AnnotateDiff(diff)
	out := FormatAnnotations(annotations)

	if !strings.Contains(out, "github.com/foo/bar") {
		t.Errorf("expected module name in output, got:\n%s", out)
	}
}
