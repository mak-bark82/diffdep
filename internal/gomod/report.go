package gomod

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ReportFormat defines the output format for a diff report.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatMarkdown ReportFormat = "markdown"
)

// WriteReport writes a human-readable diff report to the given writer.
func WriteReport(w io.Writer, diff DiffResult, format ReportFormat) error {
	switch format {
	case FormatMarkdown:
		return writeMarkdownReport(w, diff)
	default:
		return writeTextReport(w, diff)
	}
}

func writeTextReport(w io.Writer, diff DiffResult) error {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Changed) == 0 {
		_, err := fmt.Fprintln(w, "No dependency changes detected.")
		return err
	}

	keys := sortedKeys(diff.Added)
	for _, mod := range keys {
		if _, err := fmt.Fprintf(w, "+ %s %s\n", mod, diff.Added[mod]); err != nil {
			return err
		}
	}

	keys = sortedKeys(diff.Removed)
	for _, mod := range keys {
		if _, err := fmt.Fprintf(w, "- %s %s\n", mod, diff.Removed[mod]); err != nil {
			return err
		}
	}

	keys = sortedKeys(diff.Changed)
	for _, mod := range keys {
		c := diff.Changed[mod]
		if _, err := fmt.Fprintf(w, "~ %s %s -> %s\n", mod, c.From, c.To); err != nil {
			return err
		}
	}

	return nil
}

func writeMarkdownReport(w io.Writer, diff DiffResult) error {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Changed) == 0 {
		_, err := fmt.Fprintln(w, "_No dependency changes detected._")
		return err
	}

	var sb strings.Builder
	sb.WriteString("| Change | Module | From | To |\n")
	sb.WriteString("|--------|--------|------|----|\n")

	for _, mod := range sortedKeys(diff.Added) {
		sb.WriteString(fmt.Sprintf("| ✅ added | `%s` | — | `%s` |\n", mod, diff.Added[mod]))
	}
	for _, mod := range sortedKeys(diff.Removed) {
		sb.WriteString(fmt.Sprintf("| ❌ removed | `%s` | `%s` | — |\n", mod, diff.Removed[mod]))
	}
	for _, mod := range sortedKeys(diff.Changed) {
		c := diff.Changed[mod]
		sb.WriteString(fmt.Sprintf("| ⚠️ changed | `%s` | `%s` | `%s` |\n", mod, c.From, c.To))
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
