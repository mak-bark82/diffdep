package gomod

import (
	"fmt"
	"strings"
)

// DeprecationEntry holds information about a potentially deprecated dependency.
type DeprecationEntry struct {
	Module  string
	Version string
	Reason  string
	Severity string // "warn" or "info"
}

// DeprecationReport holds all detected deprecation signals.
type DeprecationReport struct {
	Branch  string
	Entries []DeprecationEntry
}

// knownDeprecations is a simple built-in registry of known deprecated modules.
var knownDeprecations = map[string]string{
	"github.com/dgrijalva/jwt-go":    "replaced by github.com/golang-jwt/jwt; unmaintained since 2021",
	"github.com/russross/blackfriday": "v1 is deprecated; migrate to v2",
	"gopkg.in/ini.v1":                 "consider migrating to a maintained config library",
	"github.com/codegangsta/cli":      "renamed to github.com/urfave/cli",
	"github.com/robfig/cron":          "v1 is deprecated; use v3",
}

// CheckDeprecations inspects a set of dependencies against known deprecated modules.
func CheckDeprecations(deps []Dependency, branch string) DeprecationReport {
	report := DeprecationReport{Branch: branch}
	for _, dep := range deps {
		if reason, ok := knownDeprecations[dep.Module]; ok {
			report.Entries = append(report.Entries, DeprecationEntry{
				Module:   dep.Module,
				Version:  dep.Version,
				Reason:   reason,
				Severity: "warn",
			})
		}
	}
	return report
}

// FormatDeprecationReport returns a human-readable report string.
func FormatDeprecationReport(r DeprecationReport) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("Deprecation check for branch %q: no known deprecated dependencies found.\n", r.Branch)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Deprecation report for branch %q:\n", r.Branch)
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  [%s] %s@%s — %s\n", strings.ToUpper(e.Severity), e.Module, e.Version, e.Reason)
	}
	return sb.String()
}
