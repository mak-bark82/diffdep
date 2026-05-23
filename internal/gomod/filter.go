package gomod

import "strings"

// FilterOptions controls which dependency changes are included in results.
type FilterOptions struct {
	// MajorOnly limits results to major version changes (e.g. v1 -> v2).
	MajorOnly bool
	// IncludeAdded includes newly added dependencies.
	IncludeAdded bool
	// IncludeRemoved includes removed dependencies.
	IncludeRemoved bool
	// PrefixFilter restricts results to modules matching the given prefix.
	PrefixFilter string
}

// FilterDiff applies the given FilterOptions to a DiffResult, returning a
// filtered copy that contains only the changes matching the criteria.
func FilterDiff(result DiffResult, opts FilterOptions) DiffResult {
	filtered := DiffResult{
		Added:   []Dependency{},
		Removed: []Dependency{},
		Changed: []DependencyChange{},
	}

	if opts.IncludeAdded {
		for _, dep := range result.Added {
			if matchesPrefix(dep.Module, opts.PrefixFilter) {
				filtered.Added = append(filtered.Added, dep)
			}
		}
	}

	if opts.IncludeRemoved {
		for _, dep := range result.Removed {
			if matchesPrefix(dep.Module, opts.PrefixFilter) {
				filtered.Removed = append(filtered.Removed, dep)
			}
		}
	}

	for _, change := range result.Changed {
		if !matchesPrefix(change.Module, opts.PrefixFilter) {
			continue
		}
		if opts.MajorOnly && !isMajorChange(change.OldVersion, change.NewVersion) {
			continue
		}
		filtered.Changed = append(filtered.Changed, change)
	}

	return filtered
}

// matchesPrefix returns true if module starts with prefix, or prefix is empty.
func matchesPrefix(module, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(module, prefix)
}

// isMajorChange returns true when the major version segment differs between
// two semver strings (e.g. "v1.2.3" vs "v2.0.0").
func isMajorChange(oldVer, newVer string) bool {
	return majorSegment(oldVer) != majorSegment(newVer)
}

// majorSegment extracts the major version prefix (e.g. "v1") from a semver.
func majorSegment(version string) string {
	v := strings.TrimPrefix(version, "v")
	if idx := strings.Index(v, "."); idx != -1 {
		return "v" + v[:idx]
	}
	return "v" + v
}
