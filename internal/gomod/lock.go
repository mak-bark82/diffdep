package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// LockEntry represents a pinned module version requirement.
type LockEntry struct {
	Module  string `json:"module"`
	Version string `json:"version"`
	Exact   bool   `json:"exact"` // if true, only this exact version is allowed
}

// LockFile holds a set of locked module versions.
type LockFile struct {
	Entries []LockEntry `json:"entries"`
}

// NewLockFile creates a LockFile from the current dependency snapshot.
func NewLockFile(deps []Dependency) *LockFile {
	entries := make([]LockEntry, 0, len(deps))
	for _, d := range deps {
		entries = append(entries, LockEntry{
			Module:  d.Module,
			Version: d.Version,
			Exact:   true,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Module < entries[j].Module
	})
	return &LockFile{Entries: entries}
}

// CheckLock verifies that all deps in the diff comply with the lock file.
// Returns a list of violation messages.
func CheckLock(lf *LockFile, diff []DiffEntry) []string {
	if lf == nil {
		return nil
	}
	lockMap := make(map[string]LockEntry, len(lf.Entries))
	for _, e := range lf.Entries {
		lockMap[e.Module] = e
	}

	var violations []string
	for _, d := range diff {
		entry, locked := lockMap[d.Module]
		if !locked {
			continue
		}
		newVer := d.NewVersion
		if newVer == "" {
			newVer = d.OldVersion
		}
		if entry.Exact && newVer != entry.Version {
			violations = append(violations, fmt.Sprintf(
				"%s: locked to %s, got %s", d.Module, entry.Version, newVer,
			))
		}
	}
	return violations
}

// FormatLockViolations formats lock violations for display.
func FormatLockViolations(violations []string) string {
	if len(violations) == 0 {
		return "No lock violations found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Lock violations (%d):\n", len(violations)))
	for _, v := range violations {
		sb.WriteString("  - " + v + "\n")
	}
	return sb.String()
}
