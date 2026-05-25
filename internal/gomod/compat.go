package gomod

import (
	"fmt"
	"strings"
)

// CompatEntry represents a compatibility check result for a single dependency.
type CompatEntry struct {
	Module  string
	From    string
	To      string
	GoMajor bool   // true if Go major version constraint is violated
	Reason  string
}

// CompatReport holds the full compatibility analysis result.
type CompatReport struct {
	Branch  string
	Entries []CompatEntry
}

// CheckCompat analyzes a diff for Go module compatibility issues.
// It flags cases where a module moves to an incompatible major version
// (e.g., v1 -> v2) without updating the import path suffix.
func CheckCompat(branch string, diff []DiffEntry) CompatReport {
	report := CompatReport{Branch: branch}
	for _, d := range diff {
		if d.Type != Changed {
			continue
		}
		if isMajorChange(d.OldVersion, d.NewVersion) && !hasMajorSuffix(d.Module, d.NewVersion) {
			report.Entries = append(report.Entries, CompatEntry{
				Module:  d.Module,
				From:    d.OldVersion,
				To:      d.NewVersion,
				GoMajor: true,
				Reason:  fmt.Sprintf("major version bump (%s -> %s) without import path suffix", d.OldVersion, d.NewVersion),
			})
		}
	}
	return report
}

// hasMajorSuffix returns true if the module path already encodes the major
// version matching the given semver string (e.g. module/v2 for v2.x.x).
func hasMajorSuffix(module, version string) bool {
	seg := majorSegment(version)
	if seg == "v1" || seg == "v0" {
		return true
	}
	return strings.HasSuffix(module, "/"+seg)
}

// FormatCompatReport returns a human-readable compatibility report.
func FormatCompatReport(r CompatReport) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("compat: no issues found on branch %q\n", r.Branch)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Compatibility issues on branch %q:\n", r.Branch)
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  [COMPAT] %s: %s\n", e.Module, e.Reason)
	}
	return sb.String()
}
