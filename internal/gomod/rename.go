package gomod

import (
	"fmt"
	"strings"
)

// RenameEntry represents a detected module rename between two dependency sets.
type RenameEntry struct {
	OldModule string
	NewModule string
	Version   string
	Note      string
}

// RenameReport holds all detected renames for a comparison.
type RenameReport struct {
	Branch  string
	Entries []RenameEntry
}

// DetectRenames compares two dependency maps and identifies likely module renames.
// A rename is inferred when a module is removed and another with a similar base path is added
// at the same or compatible version.
func DetectRenames(base, head map[string]string, branch string) RenameReport {
	report := RenameReport{Branch: branch}

	diff := DiffDependencies(DepsToMap(mapToDeps(base)), DepsToMap(mapToDeps(head)))

	for _, removed := range diff {
		if removed.OldVersion == "" {
			continue
		}
		for _, added := range diff {
			if added.NewVersion == "" {
				continue
			}
			if looksLikeRename(removed.Module, added.Module) {
				note := fmt.Sprintf("possible rename: %s -> %s", removed.Module, added.Module)
				report.Entries = append(report.Entries, RenameEntry{
					OldModule: removed.Module,
					NewModule: added.Module,
					Version:   added.NewVersion,
					Note:      note,
				})
			}
		}
	}
	return report
}

// looksLikeRename returns true if two module paths share a common base name.
func looksLikeRename(a, b string) bool {
	if a == b {
		return false
	}
	baseA := moduleBaseName(a)
	baseB := moduleBaseName(b)
	return baseA != "" && baseA == baseB
}

func moduleBaseName(mod string) string {
	parts := strings.Split(mod, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimPrefix(parts[len(parts)-1], "v")
}

func mapToDeps(m map[string]string) []Dependency {
	deps := make([]Dependency, 0, len(m))
	for mod, ver := range m {
		deps = append(deps, Dependency{Module: mod, Version: ver})
	}
	return deps
}

// FormatRenameReport formats the rename report as human-readable text.
func FormatRenameReport(r RenameReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Rename Report [branch: %s]\n", r.Branch)
	if len(r.Entries) == 0 {
		sb.WriteString("  No likely renames detected.\n")
		return sb.String()
	}
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  %s -> %s @ %s\n    note: %s\n", e.OldModule, e.NewModule, e.Version, e.Note)
	}
	return sb.String()
}
