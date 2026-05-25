package gomod

import (
	"fmt"
	"strings"
)

// LockSummary holds aggregate counts from an audit lock check.
type LockSummary struct {
	Total    int
	Locked   int
	Unlocked int
	Violations int
}

// SummarizeLock produces a LockSummary from a set of dependencies and a lock file.
func SummarizeLock(deps []Dependency, lock *LockFile) LockSummary {
	if lock == nil {
		return LockSummary{Total: len(deps), Unlocked: len(deps)}
	}
	s := LockSummary{Total: len(deps)}
	violations := CheckLock(deps, lock)
	s.Violations = len(violations)
	lockMap := make(map[string]string, len(lock.Entries))
	for _, e := range lock.Entries {
		lockMap[e.Module] = e.Version
	}
	for _, d := range deps {
		if _, ok := lockMap[d.Module]; ok {
			s.Locked++
		} else {
			s.Unlocked++
		}
	}
	return s
}

// FormatLockSummary returns a short textual summary of the lock state.
func FormatLockSummary(s LockSummary) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Lock summary: %d total, %d locked, %d unlocked, %d violations\n",
		s.Total, s.Locked, s.Unlocked, s.Violations)
	return sb.String()
}
