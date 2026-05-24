package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// PruneSuggestion represents a dependency that may be safe to remove.
type PruneSuggestion struct {
	Module  string
	Version string
	Reason  string
}

// PruneResult holds the output of a prune analysis.
type PruneResult struct {
	Suggestions []PruneSuggestion
}

// AnalyzePrune inspects the diff and current deps to suggest removals.
// It flags dependencies that were present in base but are absent in head,
// and dependencies whose versions have not changed across branches.
func AnalyzePrune(base, head []Dependency, diff []DiffEntry) PruneResult {
	headMap := DepsToMap(head)
	baseMap := DepsToMap(base)

	var suggestions []PruneSuggestion

	// Suggest removal for deps removed in head (no longer needed).
	for _, entry := range diff {
		if entry.Type == Removed {
			suggestions = append(suggestions, PruneSuggestion{
				Module:  entry.Module,
				Version: entry.OldVersion,
				Reason:  "removed in head branch; consider dropping from go.mod",
			})
		}
	}

	// Suggest review for deps unchanged across both branches (potential dead weight).
	for mod, headVer := range headMap {
		if baseVer, ok := baseMap[mod]; ok && baseVer == headVer {
			// Only suggest if not already in diff.
			if !inDiff(diff, mod) {
				suggestions = append(suggestions, PruneSuggestion{
					Module:  mod,
					Version: headVer,
					Reason:  "unchanged across branches; verify it is still required",
				})
			}
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Module < suggestions[j].Module
	})

	return PruneResult{Suggestions: suggestions}
}

func inDiff(diff []DiffEntry, module string) bool {
	for _, d := range diff {
		if d.Module == module {
			return true
		}
	}
	return false
}

// FormatPruneResult formats the prune result as human-readable text.
func FormatPruneResult(r PruneResult) string {
	if len(r.Suggestions) == 0 {
		return "No prune suggestions found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Prune Suggestions (%d):\n", len(r.Suggestions)))
	for _, s := range r.Suggestions {
		sb.WriteString(fmt.Sprintf("  - %s @ %s\n      reason: %s\n", s.Module, s.Version, s.Reason))
	}
	return sb.String()
}
