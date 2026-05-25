package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// OutdatedEntry represents a dependency that may have a newer version available.
type OutdatedEntry struct {
	Module  string
	Current string
	Latest  string
	IsStale bool
}

// OutdatedReport holds all outdated entries for a given branch.
type OutdatedReport struct {
	Branch  string
	Entries []OutdatedEntry
}

// LatestResolver is a function that resolves the latest known version for a module.
// In production this can query a registry; in tests it can be stubbed.
type LatestResolver func(module string) (string, error)

// CheckOutdated compares current deps against resolved latest versions.
func CheckOutdated(branch string, deps []Dependency, resolve LatestResolver) (*OutdatedReport, error) {
	report := &OutdatedReport{Branch: branch}
	for _, dep := range deps {
		latest, err := resolve(dep.Module)
		if err != nil {
			continue
		}
		entry := OutdatedEntry{
			Module:  dep.Module,
			Current: dep.Version,
			Latest:  latest,
			IsStale: latest != dep.Version,
		}
		report.Entries = append(report.Entries, entry)
	}
	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Module < report.Entries[j].Module
	})
	return report, nil
}

// FormatOutdatedReport formats the report as human-readable text.
func FormatOutdatedReport(r *OutdatedReport) string {
	if r == nil || len(r.Entries) == 0 {
		return fmt.Sprintf("outdated report for branch %q: all dependencies are up to date\n", r.Branch)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Outdated dependencies on branch %q:\n", r.Branch)
	for _, e := range r.Entries {
		if e.IsStale {
			fmt.Fprintf(&sb, "  %s: %s -> %s\n", e.Module, e.Current, e.Latest)
		}
	}
	staleCount := 0
	for _, e := range r.Entries {
		if e.IsStale {
			staleCount++
		}
	}
	if staleCount == 0 {
		fmt.Fprintf(&sb, "  (none)\n")
	}
	return sb.String()
}
