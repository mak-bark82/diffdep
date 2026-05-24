package gomod

import (
	"fmt"
	"strings"
)

// PinEntry represents a pinned module at a specific version.
type PinEntry struct {
	Module  string `json:"module"`
	Version string `json:"version"`
	Reason  string `json:"reason,omitempty"`
}

// PinList holds a set of pinned module versions.
type PinList struct {
	Pins []PinEntry `json:"pins"`
}

// NewPinList creates an empty PinList.
func NewPinList() *PinList {
	return &PinList{}
}

// Add appends a pin entry to the list.
func (pl *PinList) Add(module, version, reason string) {
	pl.Pins = append(pl.Pins, PinEntry{
		Module:  module,
		Version: version,
		Reason:  reason,
	})
}

// CheckViolations compares a diff against pinned versions and returns violations.
func (pl *PinList) CheckViolations(diff []DiffEntry) []PinViolation {
	pinMap := make(map[string]PinEntry, len(pl.Pins))
	for _, p := range pl.Pins {
		pinMap[p.Module] = p
	}

	var violations []PinViolation
	for _, d := range diff {
		pin, ok := pinMap[d.Module]
		if !ok {
			continue
		}
		if d.NewVersion != pin.Version {
			violations = append(violations, PinViolation{
				Module:   d.Module,
				Pinned:   pin.Version,
				Actual:   d.NewVersion,
				Reason:   pin.Reason,
			})
		}
	}
	return violations
}

// PinViolation describes a module that deviates from its pinned version.
type PinViolation struct {
	Module string
	Pinned string
	Actual string
	Reason string
}

// FormatPinViolations returns a human-readable summary of pin violations.
func FormatPinViolations(violations []PinViolation) string {
	if len(violations) == 0 {
		return "No pin violations found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pin violations (%d):\n", len(violations)))
	for _, v := range violations {
		line := fmt.Sprintf("  %s: pinned=%s actual=%s", v.Module, v.Pinned, v.Actual)
		if v.Reason != "" {
			line += fmt.Sprintf(" [%s]", v.Reason)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
