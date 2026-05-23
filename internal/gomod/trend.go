package gomod

import (
	"fmt"
	"sort"
	"time"
)

// TrendEntry represents a snapshot of dependency changes at a point in time.
type TrendEntry struct {
	Timestamp time.Time        `json:"timestamp"`
	Branch    string           `json:"branch"`
	Added     int              `json:"added"`
	Removed   int              `json:"removed"`
	Changed   int              `json:"changed"`
	Modules   []string         `json:"modules"`
}

// Trend holds a series of TrendEntry snapshots.
type Trend struct {
	Entries []TrendEntry `json:"entries"`
}

// RecordSnapshot appends a new snapshot derived from a DiffResult to the trend.
func (t *Trend) RecordSnapshot(branch string, diff []DiffEntry) {
	entry := TrendEntry{
		Timestamp: time.Now().UTC(),
		Branch:    branch,
	}
	for _, d := range diff {
		switch d.ChangeType {
		case ChangeAdded:
			entry.Added++
		case ChangeRemoved:
			entry.Removed++
		case ChangeUpdated:
			entry.Changed++
		}
		entry.Modules = append(entry.Modules, d.Module)
	}
	sort.Strings(entry.Modules)
	t.Entries = append(t.Entries, entry)
}

// Summary returns a human-readable summary of the trend.
func (t *Trend) Summary() string {
	if len(t.Entries) == 0 {
		return "no trend data recorded"
	}
	var out string
	for _, e := range t.Entries {
		out += fmt.Sprintf("[%s] branch=%s added=%d removed=%d changed=%d\n",
			e.Timestamp.Format(time.RFC3339), e.Branch, e.Added, e.Removed, e.Changed)
	}
	return out
}
