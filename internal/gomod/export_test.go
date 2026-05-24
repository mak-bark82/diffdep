package gomod

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleDiffForExport() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", Type: DiffAdded, NewVersion: "v1.2.0"},
		{Module: "github.com/baz/qux", Type: DiffChanged, OldVersion: "v1.0.0", NewVersion: "v2.0.0"},
		{Module: "github.com/old/pkg", Type: DiffRemoved, OldVersion: "v0.9.0"},
	}
}

func TestExportDiff_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := ExportDiff(&buf, sampleDiffForExport(), ExportJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
	// records should be sorted by module name
	if records[0].Module != "github.com/baz/qux" {
		t.Errorf("expected first record to be github.com/baz/qux, got %s", records[0].Module)
	}
}

func TestExportDiff_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := ExportDiff(&buf, sampleDiffForExport(), ExportCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "module") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExportDiff_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ExportDiff(&buf, sampleDiffForExport(), ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportDiff_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := ExportDiff(&buf, []DiffEntry{}, ExportJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty records, got %d", len(records))
	}
}
