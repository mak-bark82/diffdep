package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// GraphDiffEntry represents a change in the dependency graph structure,
// capturing not just the direct dependency change but also the modules
// transitively affected by it.
type GraphDiffEntry struct {
	Module    string
	ChangeType string // "added", "removed", "changed"
	OldVersion string
	NewVersion string
	Affected   []string // modules that depend on this module
}

// GraphDiff holds the full set of graph-aware dependency changes.
type GraphDiff struct {
	Entries []GraphDiffEntry
}

// DiffGraph compares two dependency graphs and produces a GraphDiff that
// annotates each changed dependency with the list of modules transitively
// affected by the change. This helps surface the blast radius of an update.
func DiffGraph(base, head *DependencyGraph, diff []DiffEntry) GraphDiff {
	var entries []GraphDiffEntry

	for _, d := range diff {
		entry := GraphDiffEntry{
			Module:     d.Module,
			ChangeType: d.ChangeType,
			OldVersion: d.OldVersion,
			NewVersion: d.NewVersion,
		}

		// Use the head graph to find affected modules when available;
		// fall back to the base graph for removals.
		var g *DependencyGraph
		if head != nil {
			g = head
		} else {
			g = base
		}

		if g != nil {
			affected := g.Affected(d.Module)
			sort.Strings(affected)
			entry.Affected = affected
		}

		entries = append(entries, entry)
	}

	// Sort entries by module name for deterministic output.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Module < entries[j].Module
	})

	return GraphDiff{Entries: entries}
}

// FormatGraphDiff formats a GraphDiff as a human-readable text report.
// Each changed module is listed with its version change and the set of
// modules that will be affected by the change.
func FormatGraphDiff(gd GraphDiff) string {
	if len(gd.Entries) == 0 {
		return "No dependency graph changes detected.\n"
	}

	var sb strings.Builder
	sb.WriteString("Dependency Graph Diff\n")
	sb.WriteString(strings.Repeat("=", 40) + "\n")

	for _, e := range gd.Entries {
		switch e.ChangeType {
		case "added":
			sb.WriteString(fmt.Sprintf("+ %s %s (new)\n", e.Module, e.NewVersion))
		case "removed":
			sb.WriteString(fmt.Sprintf("- %s %s (removed)\n", e.Module, e.OldVersion))
		case "changed":
			sb.WriteString(fmt.Sprintf("~ %s %s -> %s\n", e.Module, e.OldVersion, e.NewVersion))
		default:
			sb.WriteString(fmt.Sprintf("? %s\n", e.Module))
		}

		if len(e.Affected) > 0 {
			sb.WriteString(fmt.Sprintf("  Affects (%d): %s\n", len(e.Affected), strings.Join(e.Affected, ", ")))
		} else {
			sb.WriteString("  Affects: (none)\n")
		}
	}

	return sb.String()
}
