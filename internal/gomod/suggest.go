package gomod

import (
	"fmt"
	"strings"
)

// Suggestion represents an actionable recommendation for a dependency change.
type Suggestion struct {
	Module  string
	Current string
	Target  string
	Reason  string
	Priority string // "high", "medium", "low"
}

// String returns a human-readable representation of the suggestion.
func (s Suggestion) String() string {
	return fmt.Sprintf("[%s] %s: %s -> %s (%s)", strings.ToUpper(s.Priority), s.Module, s.Current, s.Target, s.Reason)
}

// GenerateSuggestions inspects a diff and produces actionable upgrade or
// downgrade suggestions based on version change patterns.
func GenerateSuggestions(diff []DiffEntry) []Suggestion {
	var suggestions []Suggestion

	for _, entry := range diff {
		switch entry.Change {
		case ChangeAdded:
			suggestions = append(suggestions, Suggestion{
				Module:   entry.Module,
				Current:  "(none)",
				Target:   entry.NewVersion,
				Reason:   "new dependency introduced; verify it is intentional",
				Priority: "low",
			})

		case ChangeRemoved:
			suggestions = append(suggestions, Suggestion{
				Module:   entry.Module,
				Current:  entry.OldVersion,
				Target:   "(removed)",
				Reason:   "dependency removed; ensure no transitive consumers remain",
				Priority: "medium",
			})

		case ChangeUpdated:
			priority := "low"
			reason := "minor or patch version bump"
			if isMajorChange(entry.OldVersion, entry.NewVersion) {
				priority = "high"
				reason = "major version bump — review breaking changes in upstream changelog"
			} else if isDowngrade(entry.OldVersion, entry.NewVersion) {
				priority = "medium"
				reason = "version downgrade detected; confirm this is deliberate"
			}
			suggestions = append(suggestions, Suggestion{
				Module:   entry.Module,
				Current:  entry.OldVersion,
				Target:   entry.NewVersion,
				Reason:   reason,
				Priority: priority,
			})
		}
	}

	return suggestions
}

// isDowngrade returns true when newVer sorts lexicographically before oldVer.
// This is a best-effort heuristic for semver strings like "v1.2.3".
func isDowngrade(oldVer, newVer string) bool {
	return newVer < oldVer
}

// FormatSuggestions renders suggestions as a plain-text report.
func FormatSuggestions(suggestions []Suggestion) string {
	if len(suggestions) == 0 {
		return "No suggestions.\n"
	}
	var sb strings.Builder
	sb.WriteString("Suggestions:\n")
	for _, s := range suggestions {
		sb.WriteString("  " + s.String() + "\n")
	}
	return sb.String()
}
