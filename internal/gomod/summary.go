package gomod

import "fmt"

// Summary holds aggregate statistics for a dependency diff.
type Summary struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// Summarize computes a Summary from a slice of DiffEntries.
func Summarize(entries []DiffEntry) Summary {
	s := Summary{Total: len(entries)}
	for _, e := range entries {
		switch e.ChangeType {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Changed:
			s.Changed++
		}
	}
	return s
}

// String returns a human-readable one-line summary.
func (s Summary) String() string {
	return fmt.Sprintf("total=%d added=%d removed=%d changed=%d",
		s.Total, s.Added, s.Removed, s.Changed)
}

// HasChanges returns true when any dependency change exists.
func (s Summary) HasChanges() bool {
	return s.Total > 0
}
