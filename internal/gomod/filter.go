package gomod

import "strings"

// FilterOptions controls which diff entries are included in output.
type FilterOptions struct {
	MajorOnly      bool
	Prefix         string
	IncludeAdded   bool
	IncludeRemoved bool
}

// FilterDiff applies the given options to narrow down a list of DiffEntries.
func FilterDiff(entries []DiffEntry, opts FilterOptions) []DiffEntry {
	var out []DiffEntry
	for _, e := range entries {
		if opts.Prefix != "" && !matchesPrefix(e.Module, opts.Prefix) {
			continue
		}
		switch e.ChangeType {
		case Added:
			if !opts.IncludeAdded {
				continue
			}
		case Removed:
			if !opts.IncludeRemoved {
				continue
			}
		case Changed:
			if opts.MajorOnly && !isMajorChange(e.OldVersion, e.NewVersion) {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}

func matchesPrefix(module, prefix string) bool {
	return strings.HasPrefix(module, prefix)
}

func isMajorChange(oldVer, newVer string) bool {
	return majorSegment(oldVer) != majorSegment(newVer)
}

func majorSegment(version string) string {
	v := strings.TrimPrefix(version, "v")
	parts := strings.SplitN(v, ".", 2)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
