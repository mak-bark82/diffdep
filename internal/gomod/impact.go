package gomod

import (
	"fmt"
	"strings"
)

// ImpactLevel describes how broadly a dependency change may affect a project.
type ImpactLevel string

const (
	ImpactLow      ImpactLevel = "low"
	ImpactMedium   ImpactLevel = "medium"
	ImpactHigh     ImpactLevel = "high"
	ImpactCritical ImpactLevel = "critical"
)

// ImpactEntry records the assessed impact for a single dependency change.
type ImpactEntry struct {
	Module  string
	OldVer  string
	NewVer  string
	Kind    ChangeKind
	Level   ImpactLevel
	Reason  string
}

// ImpactReport holds all entries produced by AssessImpact.
type ImpactReport struct {
	Branch  string
	Entries []ImpactEntry
}

// AssessImpact evaluates each diff entry and assigns an impact level.
func AssessImpact(branch string, diff []DiffEntry) ImpactReport {
	report := ImpactReport{Branch: branch}
	for _, d := range diff {
		entry := ImpactEntry{
			Module: d.Module,
			OldVer: d.OldVersion,
			NewVer: d.NewVersion,
			Kind:   d.Kind,
		}
		switch d.Kind {
		case KindAdded:
			entry.Level = ImpactLow
			entry.Reason = "new dependency introduced"
		case KindRemoved:
			entry.Level = ImpactMedium
			entry.Reason = "dependency removed; consumers may break"
		case KindChanged:
			if isMajorChange(d.OldVersion, d.NewVersion) {
				entry.Level = ImpactCritical
				entry.Reason = "major version bump; breaking API changes likely"
			} else if isDowngrade(d.OldVersion, d.NewVersion) {
				entry.Level = ImpactHigh
				entry.Reason = "version downgrade detected"
			} else {
				entry.Level = ImpactLow
				entry.Reason = "minor or patch update"
			}
		}
		report.Entries = append(report.Entries, entry)
	}
	return report
}

// FormatImpactReport returns a human-readable summary of the impact report.
func FormatImpactReport(r ImpactReport) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("Impact report for branch %q: no changes detected.\n", r.Branch)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Impact report for branch %q:\n", r.Branch)
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  [%s] %s %s -> %s — %s\n",
			strings.ToUpper(string(e.Level)), e.Module, e.OldVer, e.NewVer, e.Reason)
	}
	return sb.String()
}
