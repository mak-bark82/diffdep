package gomod

import (
	"fmt"
	"strings"
)

// LockSummary holds aggregated statistics about a lock check.
type LockSummary struct {
	TotalLocked    int
	TotalChecked   int
	ViolationCount int
	Violations     []string
}

// SummarizeLock produces a LockSummary from a LockFile and a set of violations.
func SummarizeLock(lf *LockFile, diff []DiffEntry, violations []string) LockSummary {
	locked := 0
	if lf != nil {
		locked = len(lf.Entries)
	}
	return LockSummary{
		TotalLocked:    locked,
		TotalChecked:   len(diff),
		ViolationCount: len(violations),
		Violations:     violations,
	}
}

// String returns a human-readable summary of the lock check.
func (s LockSummary) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Lock summary: %d locked, %d checked, %d violation(s)\n",
		s.TotalLocked, s.TotalChecked, s.ViolationCount))
	for _, v := range s.Violations {
		sb.WriteString("  ! " + v + "\n")
	}
	return sb.String()
}
