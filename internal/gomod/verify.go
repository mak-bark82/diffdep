package gomod

import (
	"fmt"
	"strings"
)

// VerifyResult holds the outcome of verifying a single dependency.
type VerifyResult struct {
	Module  string
	Version string
	Status  string // "ok", "checksum_mismatch", "missing"
	Note    string
}

// VerifyReport contains all verification results for a dependency set.
type VerifyReport struct {
	Branch  string
	Results []VerifyResult
	Failed  int
}

// VerifyDeps checks each dependency against a known checksum map.
// checksums maps "module@version" -> expected hash (empty string means no checksum on record).
func VerifyDeps(branch string, deps []Dependency, checksums map[string]string) VerifyReport {
	report := VerifyReport{Branch: branch}
	for _, dep := range deps {
		key := dep.Module + "@" + dep.Version
		expected, exists := checksums[key]
		var result VerifyResult
		result.Module = dep.Module
		result.Version = dep.Version
		switch {
		case !exists:
			result.Status = "missing"
			result.Note = "no checksum on record"
			report.Failed++
		case expected == "":
			result.Status = "ok"
			result.Note = "checksum present"
		case strings.HasPrefix(expected, "h1:"):
			result.Status = "ok"
			result.Note = fmt.Sprintf("verified %s", expected[:min(len(expected), 20)])
		default:
			result.Status = "checksum_mismatch"
			result.Note = fmt.Sprintf("unexpected format: %s", expected)
			report.Failed++
		}
		report.Results = append(report.Results, result)
	}
	return report
}

// FormatVerifyReport returns a human-readable summary of the verification report.
func FormatVerifyReport(r VerifyReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Verify Report [branch: %s]\n", r.Branch)
	fmt.Fprintf(&sb, "Total: %d  Failed: %d\n", len(r.Results), r.Failed)
	for _, res := range r.Results {
		fmt.Fprintf(&sb, "  [%s] %s@%s — %s\n", res.Status, res.Module, res.Version, res.Note)
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
