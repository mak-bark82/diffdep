package gomod

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ChangelogEntry represents a single versioned changelog record.
type ChangelogEntry struct {
	Timestamp time.Time
	Branch    string
	Diff      []DiffEntry
	Summary   Summary
}

// Changelog holds an ordered list of entries.
type Changelog struct {
	Entries []ChangelogEntry
}

// AddEntry appends a new entry to the changelog.
func (c *Changelog) AddEntry(branch string, diff []DiffEntry) {
	entry := ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Branch:    branch,
		Diff:      diff,
		Summary:   Summarize(diff),
	}
	c.Entries = append(c.Entries, entry)
}

// FormatText renders the changelog as plain text.
func (c *Changelog) FormatText() string {
	var sb strings.Builder
	for _, entry := range c.Entries {
		sb.WriteString(fmt.Sprintf("[%s] branch=%s %s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Branch,
			entry.Summary.String(),
		))
		keys := make([]string, 0, len(entry.Diff))
		for _, d := range entry.Diff {
			keys = append(keys, d.Module)
		}
		sort.Strings(keys)
		modMap := make(map[string]DiffEntry)
		for _, d := range entry.Diff {
			modMap[d.Module] = d
		}
		for _, k := range keys {
			d := modMap[k]
			sb.WriteString(fmt.Sprintf("  %s %s -> %s (%s)\n", d.Module, d.OldVersion, d.NewVersion, d.ChangeType))
		}
	}
	return sb.String()
}

// Len returns the number of entries in the changelog.
func (c *Changelog) Len() int {
	return len(c.Entries)
}
