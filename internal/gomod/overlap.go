package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// OverlapEntry represents a module present in both branches with version details.
type OverlapEntry struct {
	Module  string
	BaseVer string
	HeadVer string
	Same    bool
}

// OverlapReport holds the result of comparing two dependency sets.
type OverlapReport struct {
	Branch  string
	Entries []OverlapEntry
	Shared  int
	Diverge int
}

// AnalyzeOverlap identifies modules present in both base and head dependency
// sets and reports whether their versions agree or diverge.
func AnalyzeOverlap(branch string, base, head []Dependency) OverlapReport {
	baseMap := DepsToMap(base)
	headMap := DepsToMap(head)

	seen := make(map[string]bool)
	var entries []OverlapEntry

	for mod, bVer := range baseMap {
		if hVer, ok := headMap[mod]; ok {
			same := bVer == hVer
			entries = append(entries, OverlapEntry{
				Module:  mod,
				BaseVer: bVer,
				HeadVer: hVer,
				Same:    same,
			})
			seen[mod] = true
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Module < entries[j].Module
	})

	shared, diverge := 0, 0
	for _, e := range entries {
		if e.Same {
			shared++
		} else {
			diverge++
		}
	}

	_ = seen
	return OverlapReport{
		Branch:  branch,
		Entries: entries,
		Shared:  shared,
		Diverge: diverge,
	}
}

// FormatOverlapReport returns a human-readable summary of the overlap report.
func FormatOverlapReport(r OverlapReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Overlap Report [branch: %s]\n", r.Branch)
	fmt.Fprintf(&sb, "  Shared (same version): %d\n", r.Shared)
	fmt.Fprintf(&sb, "  Diverged (version differs): %d\n", r.Diverge)
	if len(r.Entries) == 0 {
		sb.WriteString("  No shared modules found.\n")
		return sb.String()
	}
	sb.WriteString("\n")
	for _, e := range r.Entries {
		if e.Same {
			fmt.Fprintf(&sb, "  [=] %s @ %s\n", e.Module, e.BaseVer)
		} else {
			fmt.Fprintf(&sb, "  [~] %s  base=%s  head=%s\n", e.Module, e.BaseVer, e.HeadVer)
		}
	}
	return sb.String()
}
