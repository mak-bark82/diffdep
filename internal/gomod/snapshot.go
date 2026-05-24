package gomod

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Snapshot represents the state of dependencies at a point in time.
type Snapshot struct {
	Branch    string            `json:"branch"`
	Timestamp time.Time         `json:"timestamp"`
	Deps      []Dependency      `json:"deps"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// NewSnapshot creates a Snapshot from a list of dependencies.
func NewSnapshot(branch string, deps []Dependency) Snapshot {
	sorted := make([]Dependency, len(deps))
	copy(sorted, deps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Module < sorted[j].Module
	})
	return Snapshot{
		Branch:    branch,
		Timestamp: time.Now().UTC(),
		Deps:      sorted,
	}
}

// DiffSnapshot computes the DiffResult between two snapshots.
func DiffSnapshot(base, head Snapshot) DiffResult {
	baseMap := DepsToMap(base.Deps)
	headMap := DepsToMap(head.Deps)
	return DiffDependencies(baseMap, headMap)
}

// SnapshotSummary returns a human-readable summary line for a snapshot.
func SnapshotSummary(s Snapshot) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "branch=%s time=%s deps=%d",
		s.Branch,
		s.Timestamp.Format(time.RFC3339),
		len(s.Deps),
	)
	return sb.String()
}
