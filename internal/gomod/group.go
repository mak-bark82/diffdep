package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// DiffGroup represents a named collection of dependency changes.
type DiffGroup struct {
	Name    string
	Entries []DiffEntry
}

// GroupByOrg groups diff entries by their Go module organization/host prefix.
// For example, "github.com/foo/bar" and "github.com/foo/baz" both fall under "github.com/foo".
func GroupByOrg(diff []DiffEntry) []DiffGroup {
	groupMap := make(map[string][]DiffEntry)

	for _, entry := range diff {
		org := orgPrefix(entry.Module)
		groupMap[org] = append(groupMap[org], entry)
	}

	groups := make([]DiffGroup, 0, len(groupMap))
	for name, entries := range groupMap {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Module < entries[j].Module
		})
		groups = append(groups, DiffGroup{Name: name, Entries: entries})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups
}

// FormatGroupReport returns a human-readable grouped diff report.
func FormatGroupReport(groups []DiffGroup) string {
	if len(groups) == 0 {
		return "No dependency changes detected.\n"
	}

	var sb strings.Builder
	for _, g := range groups {
		fmt.Fprintf(&sb, "[%s] (%d change(s))\n", g.Name, len(g.Entries))
		for _, e := range g.Entries {
			switch e.ChangeType {
			case ChangeAdded:
				fmt.Fprintf(&sb, "  + %s %s\n", e.Module, e.NewVersion)
			case ChangeRemoved:
				fmt.Fprintf(&sb, "  - %s %s\n", e.Module, e.OldVersion)
			case ChangeUpdated:
				fmt.Fprintf(&sb, "  ~ %s %s -> %s\n", e.Module, e.OldVersion, e.NewVersion)
			}
		}
	}
	return sb.String()
}

// orgPrefix extracts the org/host prefix from a module path.
// e.g. "github.com/foo/bar" -> "github.com/foo"
// e.g. "golang.org/x/text" -> "golang.org/x"
func orgPrefix(module string) string {
	parts := strings.Split(module, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[:2], "/")
	}
	return parts[0]
}
