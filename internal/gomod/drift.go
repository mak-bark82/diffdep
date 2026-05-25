package gomod

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// DriftEntry describes how far a dependency has drifted from a reference point.
type DriftEntry struct {
	Module   string
	Baseline string
	Current  string
	Days     int
	Severity string // "low", "medium", "high"
}

// DriftReport holds all drift entries for a given branch.
type DriftReport struct {
	Branch    string
	CreatedAt time.Time
	Entries   []DriftEntry
}

// AnalyzeDrift compares current deps against a baseline snapshot and computes drift.
func AnalyzeDrift(branch string, baseline, current []Dependency, since time.Time) DriftReport {
	report := DriftReport{
		Branch:    branch,
		CreatedAt: time.Now(),
	}

	baseMap := DepsToMap(baseline)
	currMap := DepsToMap(current)
	daysSince := int(time.Since(since).Hours() / 24)

	for mod, currVer := range currMap {
		baseVer, existed := baseMap[mod]
		if !existed {
			continue
		}
		if baseVer == currVer {
			continue
		}
		entry := DriftEntry{
			Module:   mod,
			Baseline: baseVer,
			Current:  currVer,
			Days:     daysSince,
			Severity: driftSeverity(baseVer, currVer, daysSince),
		}
		report.Entries = append(report.Entries, entry)
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Module < report.Entries[j].Module
	})
	return report
}

func driftSeverity(baseline, current string, days int) string {
	if isMajorChange(baseline, current) {
		return "high"
	}
	if days > 90 {
		return "medium"
	}
	return "low"
}

// FormatDriftReport returns a human-readable drift report.
func FormatDriftReport(r DriftReport) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("Drift report for %s: no drift detected.\n", r.Branch)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Drift report for branch %q (as of %s):\n", r.Branch, r.CreatedAt.Format("2006-01-02"))
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  [%s] %s: %s -> %s (%d days)\n", e.Severity, e.Module, e.Baseline, e.Current, e.Days)
	}
	return sb.String()
}
