package gomod

import (
	"fmt"
	"strings"
	"time"
)

// AuditEntry represents a single audit record for a dependency change.
type AuditEntry struct {
	Module    string
	OldVersion string
	NewVersion string
	ChangeType string // added, removed, changed
	Risk       string
	Timestamp  time.Time
	Branch     string
}

// AuditReport holds all audit entries for a diff.
type AuditReport struct {
	Branch    string
	CreatedAt time.Time
	Entries   []AuditEntry
}

// NewAuditReport builds an AuditReport from a DiffResult.
func NewAuditReport(branch string, diff DiffResult) AuditReport {
	report := AuditReport{
		Branch:    branch,
		CreatedAt: time.Now().UTC(),
	}

	for _, dep := range diff.Added {
		report.Entries = append(report.Entries, AuditEntry{
			Module:     dep.Module,
			NewVersion: dep.Version,
			ChangeType: "added",
			Risk:       "low",
			Timestamp:  report.CreatedAt,
			Branch:     branch,
		})
	}

	for _, dep := range diff.Removed {
		report.Entries = append(report.Entries, AuditEntry{
			Module:     dep.Module,
			OldVersion: dep.Version,
			ChangeType: "removed",
			Risk:       "medium",
			Timestamp:  report.CreatedAt,
			Branch:     branch,
		})
	}

	for _, ch := range diff.Changed {
		risk := "low"
		if isMajorChange(ch.OldVersion, ch.NewVersion) {
			risk = "high"
		}
		report.Entries = append(report.Entries, AuditEntry{
			Module:     ch.Module,
			OldVersion: ch.OldVersion,
			NewVersion: ch.NewVersion,
			ChangeType: "changed",
			Risk:       risk,
			Timestamp:  report.CreatedAt,
			Branch:     branch,
		})
	}

	return report
}

// FormatAuditReport returns a human-readable audit report string.
func FormatAuditReport(r AuditReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Audit Report — Branch: %s\n", r.Branch)
	fmt.Fprintf(&sb, "Generated: %s\n", r.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "Entries: %d\n\n", len(r.Entries))
	for _, e := range r.Entries {
		switch e.ChangeType {
		case "added":
			fmt.Fprintf(&sb, "  [ADDED]   %s %s (risk: %s)\n", e.Module, e.NewVersion, e.Risk)
		case "removed":
			fmt.Fprintf(&sb, "  [REMOVED] %s %s (risk: %s)\n", e.Module, e.OldVersion, e.Risk)
		case "changed":
			fmt.Fprintf(&sb, "  [CHANGED] %s %s -> %s (risk: %s)\n", e.Module, e.OldVersion, e.NewVersion, e.Risk)
		}
	}
	return sb.String()
}
