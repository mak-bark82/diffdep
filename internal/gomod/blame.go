package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// BlameEntry associates a dependency change with a branch and optional note.
type BlameEntry struct {
	Module  string
	OldVer  string
	NewVer  string
	Kind    string // added, removed, changed
	Branch  string
	Note    string
}

// BlameReport holds all blame entries for a diff.
type BlameReport struct {
	Branch  string
	Entries []BlameEntry
}

// NewBlameReport builds a BlameReport from a DiffResult, attributing changes to branch.
func NewBlameReport(branch string, diff DiffResult) BlameReport {
	report := BlameReport{Branch: branch}

	for _, d := range diff.Added {
		report.Entries = append(report.Entries, BlameEntry{
			Module: d.Module,
			NewVer: d.Version,
			Kind:   "added",
			Branch: branch,
			Note:   "new dependency introduced",
		})
	}
	for _, d := range diff.Removed {
		report.Entries = append(report.Entries, BlameEntry{
			Module: d.Module,
			OldVer: d.Version,
			Kind:   "removed",
			Branch: branch,
			Note:   "dependency dropped",
		})
	}
	for _, d := range diff.Changed {
		note := "version bump"
		if isMajorChange(d.OldVersion, d.NewVersion) {
			note = "major version change"
		} else if isDowngrade(d.OldVersion, d.NewVersion) {
			note = "downgrade detected"
		}
		report.Entries = append(report.Entries, BlameEntry{
			Module: d.Module,
			OldVer: d.OldVersion,
			NewVer: d.NewVersion,
			Kind:   "changed",
			Branch: branch,
			Note:   note,
		})
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Module < report.Entries[j].Module
	})
	return report
}

// FormatBlameReport returns a human-readable blame report.
func FormatBlameReport(r BlameReport) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("blame: no changes detected on branch %q\n", r.Branch)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Blame Report — branch: %s\n", r.Branch)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 60))
	for _, e := range r.Entries {
		switch e.Kind {
		case "added":
			fmt.Fprintf(&sb, "[added]   %s @ %s — %s\n", e.Module, e.NewVer, e.Note)
		case "removed":
			fmt.Fprintf(&sb, "[removed] %s @ %s — %s\n", e.Module, e.OldVer, e.Note)
		case "changed":
			fmt.Fprintf(&sb, "[changed] %s: %s → %s — %s\n", e.Module, e.OldVer, e.NewVer, e.Note)
		}
	}
	return sb.String()
}
