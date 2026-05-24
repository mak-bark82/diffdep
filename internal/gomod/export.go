package gomod

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ExportFormat defines the output format for exporting diff results.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)

// ExportRecord represents a single row in an export.
type ExportRecord struct {
	Module  string `json:"module"`
	Change  string `json:"change"`
	OldVersion string `json:"old_version,omitempty"`
	NewVersion string `json:"new_version,omitempty"`
}

// ExportDiff writes the diff result to w in the specified format.
func ExportDiff(w io.Writer, diff []DiffEntry, format ExportFormat) error {
	records := buildRecords(diff)
	switch format {
	case ExportJSON:
		return exportJSON(w, records)
	case ExportCSV:
		return exportCSV(w, records)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func buildRecords(diff []DiffEntry) []ExportRecord {
	sorted := make([]DiffEntry, len(diff))
	copy(sorted, diff)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Module < sorted[j].Module
	})
	records := make([]ExportRecord, 0, len(sorted))
	for _, d := range sorted {
		records = append(records, ExportRecord{
			Module:     d.Module,
			Change:     string(d.Type),
			OldVersion: d.OldVersion,
			NewVersion: d.NewVersion,
		})
	}
	return records
}

func exportJSON(w io.Writer, records []ExportRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func exportCSV(w io.Writer, records []ExportRecord) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"module", "change", "old_version", "new_version"}); err != nil {
		return err
	}
	for _, r := range records {
		if err := cw.Write([]string{r.Module, r.Change, r.OldVersion, r.NewVersion}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
