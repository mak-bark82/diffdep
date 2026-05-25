package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// LicenseRisk represents the risk level of a license change.
type LicenseRisk string

const (
	LicenseRiskLow      LicenseRisk = "low"
	LicenseRiskMedium   LicenseRisk = "medium"
	LicenseRiskHigh     LicenseRisk = "high"
	LicenseRiskUnknown  LicenseRisk = "unknown"
)

// KnownRestrictiveLicenses lists license identifiers considered high-risk
// in a commercial or proprietary context.
var KnownRestrictiveLicenses = []string{
	"GPL-2.0", "GPL-3.0", "AGPL-3.0", "LGPL-2.1", "LGPL-3.0",
	"SSPL-1.0", "BUSL-1.1", "Commons-Clause",
}

// KnownPermissiveLicenses lists license identifiers generally considered
// low-risk for most projects.
var KnownPermissiveLicenses = []string{
	"MIT", "Apache-2.0", "BSD-2-Clause", "BSD-3-Clause",
	"ISC", "0BSD", "Unlicense", "CC0-1.0",
}

// LicenseEntry records a module and its associated license information.
type LicenseEntry struct {
	Module  string
	Version string
	License string
	Risk    LicenseRisk
}

// LicenseReport holds the result of a license check over a dependency diff.
type LicenseReport struct {
	Entries  []LicenseEntry
	HighRisk []LicenseEntry
}

// licenseRiskFor classifies a license string into a LicenseRisk level.
func licenseRiskFor(license string) LicenseRisk {
	if license == "" {
		return LicenseRiskUnknown
	}
	upper := strings.ToUpper(license)
	for _, r := range KnownRestrictiveLicenses {
		if strings.ToUpper(r) == upper {
			return LicenseRiskHigh
		}
	}
	for _, p := range KnownPermissiveLicenses {
		if strings.ToUpper(p) == upper {
			return LicenseRiskLow
		}
	}
	return LicenseRiskMedium
}

// CheckLicenses evaluates the license metadata supplied in the licenseMap
// against the added or changed dependencies in diff. The licenseMap should
// map module paths to SPDX license identifiers (e.g. "MIT", "Apache-2.0").
// Dependencies not present in licenseMap are recorded as unknown.
func CheckLicenses(diff []DiffEntry, licenseMap map[string]string) LicenseReport {
	var report LicenseReport

	for _, entry := range diff {
		// Only inspect newly introduced or changed dependencies.
		if entry.Type == DiffRemoved {
			continue
		}
		license := licenseMap[entry.Module]
		risk := licenseRiskFor(license)
		le := LicenseEntry{
			Module:  entry.Module,
			Version: entry.NewVersion,
			License: license,
			Risk:    risk,
		}
		report.Entries = append(report.Entries, le)
		if risk == LicenseRiskHigh || risk == LicenseRiskUnknown {
			report.HighRisk = append(report.HighRisk, le)
		}
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Module < report.Entries[j].Module
	})
	sort.Slice(report.HighRisk, func(i, j int) bool {
		return report.HighRisk[i].Module < report.HighRisk[j].Module
	})
	return report
}

// FormatLicenseReport returns a human-readable summary of the license report.
func FormatLicenseReport(r LicenseReport) string {
	if len(r.Entries) == 0 {
		return "No new or changed dependencies to check.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("License Report (%d dependencies):\n", len(r.Entries)))
	for _, e := range r.Entries {
		license := e.License
		if license == "" {
			license = "(unknown)"
		}
		sb.WriteString(fmt.Sprintf("  %-45s %-15s [%s]\n", e.Module+"@"+e.Version, license, e.Risk))
	}
	if len(r.HighRisk) > 0 {
		sb.WriteString(fmt.Sprintf("\n⚠  %d high-risk or unknown license(s) detected.\n", len(r.HighRisk)))
	}
	return sb.String()
}
