package gomod

import (
	"strings"
	"testing"
)

func TestWriteReport_NoChanges(t *testing.T) {
	diff := DiffResult{}
	var sb strings.Builder
	if err := WriteReport(&sb, diff, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No dependency changes") {
		t.Errorf("expected no-change message, got: %q", sb.String())
	}
}

func TestWriteReport_TextFormat(t *testing.T) {
	diff := DiffResult{
		Added:   map[string]string{"github.com/new/pkg": "v1.2.0"},
		Removed: map[string]string{"github.com/old/pkg": "v0.9.0"},
		Changed: map[string]VersionChange{
			"github.com/foo/bar": {From: "v1.0.0", To: "v2.0.0"},
		},
	}

	var sb strings.Builder
	if err := WriteReport(&sb, diff, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()

	if !strings.Contains(out, "+ github.com/new/pkg v1.2.0") {
		t.Errorf("missing added line in output: %q", out)
	}
	if !strings.Contains(out, "- github.com/old/pkg v0.9.0") {
		t.Errorf("missing removed line in output: %q", out)
	}
	if !strings.Contains(out, "~ github.com/foo/bar v1.0.0 -> v2.0.0") {
		t.Errorf("missing changed line in output: %q", out)
	}
}

func TestWriteReport_MarkdownFormat(t *testing.T) {
	diff := DiffResult{
		Added:   map[string]string{"github.com/new/pkg": "v1.2.0"},
		Removed: map[string]string{},
		Changed: map[string]VersionChange{
			"github.com/foo/bar": {From: "v1.0.0", To: "v2.0.0"},
		},
	}

	var sb strings.Builder
	if err := WriteReport(&sb, diff, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()

	if !strings.Contains(out, "| Change | Module | From | To |") {
		t.Errorf("missing markdown table header: %q", out)
	}
	if !strings.Contains(out, "`github.com/new/pkg`") {
		t.Errorf("missing added module in markdown: %q", out)
	}
	if !strings.Contains(out, "`v2.0.0`") {
		t.Errorf("missing target version in markdown: %q", out)
	}
}

func TestWriteReport_MarkdownNoChanges(t *testing.T) {
	diff := DiffResult{}
	var sb strings.Builder
	if err := WriteReport(&sb, diff, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "_No dependency changes detected._") {
		t.Errorf("expected markdown no-change message")
	}
}
