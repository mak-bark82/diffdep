package gomod

import (
	"fmt"
	"strings"
)

// RetractedModule represents a module version that has been retracted.
type RetractedModule struct {
	Module  string
	Version string
	Reason  string
}

// RetractReport holds the result of scanning a diff for retracted versions.
type RetractReport struct {
	Branch    string
	Retracted []RetractedModule
}

// knownRetractions is a simple in-memory registry of known retracted versions.
// In a real implementation this could be loaded from a file or remote source.
var knownRetractions = map[string]string{
	"github.com/pkg/errors@v0.8.0":           "use v0.9.1+",
	"github.com/gogo/protobuf@v1.3.1":        "security vulnerability CVE-2021-3121",
	"github.com/dgrijalva/jwt-go@v3.2.0+incompatible": "use github.com/golang-jwt/jwt instead",
}

// CheckRetractions scans the given dependency list for known retracted versions.
func CheckRetractions(branch string, deps []Dependency) RetractReport {
	report := RetractReport{Branch: branch}
	for _, dep := range deps {
		key := dep.Module + "@" + dep.Version
		if reason, ok := knownRetractions[key]; ok {
			report.Retracted = append(report.Retracted, RetractedModule{
				Module:  dep.Module,
				Version: dep.Version,
				Reason:  reason,
			})
		}
	}
	return report
}

// FormatRetractReport formats the retraction report as human-readable text.
func FormatRetractReport(r RetractReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Retraction Report [branch: %s]\n", r.Branch))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	if len(r.Retracted) == 0 {
		sb.WriteString("No retracted modules detected.\n")
		return sb.String()
	}
	for _, m := range r.Retracted {
		sb.WriteString(fmt.Sprintf("  RETRACTED %s@%s\n", m.Module, m.Version))
		sb.WriteString(fmt.Sprintf("    Reason: %s\n", m.Reason))
	}
	return sb.String()
}
