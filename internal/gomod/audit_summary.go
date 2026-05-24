package gomod

import (
	"fmt"
	"strings"
)

// AuditSummary holds aggregated counts from an AuditReport.
type AuditSummary struct {
	Branch       string
	TotalEntries int
	HighRisk     int
	MediumRisk   int
	LowRisk      int
	Added        int
	Removed      int
	Changed      int
}

// SummarizeAudit computes an AuditSummary from an AuditReport.
func SummarizeAudit(r AuditReport) AuditSummary {
	s := AuditSummary{
		Branch:       r.Branch,
		TotalEntries: len(r.Entries),
	}
	for _, e := range r.Entries {
		switch e.Risk {
		case "high":
			s.HighRisk++
		case "medium":
			s.MediumRisk++
		default:
			s.LowRisk++
		}
		switch e.ChangeType {
		case "added":
			s.Added++
		case "removed":
			s.Removed++
		case "changed":
			s.Changed++
		}
	}
	return s
}

// String returns a compact summary string.
func (s AuditSummary) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Branch: %s | Total: %d", s.Branch, s.TotalEntries)
	fmt.Fprintf(&sb, " | Added: %d Removed: %d Changed: %d", s.Added, s.Removed, s.Changed)
	fmt.Fprintf(&sb, " | Risk — High: %d Medium: %d Low: %d", s.HighRisk, s.MediumRisk, s.LowRisk)
	return sb.String()
}
