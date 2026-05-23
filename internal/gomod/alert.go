package gomod

import (
	"fmt"
	"strings"
)

// AlertLevel represents the severity of a dependency alert.
type AlertLevel string

const (
	AlertInfo     AlertLevel = "INFO"
	AlertWarning  AlertLevel = "WARNING"
	AlertCritical AlertLevel = "CRITICAL"
)

// Alert represents a single dependency change alert.
type Alert struct {
	Module  string
	Level   AlertLevel
	Message string
}

// AlertConfig controls which changes trigger alerts.
type AlertConfig struct {
	OnAdded   bool
	OnRemoved bool
	OnMajor   bool
	OnMinor   bool
}

// DefaultAlertConfig returns a config that alerts on major changes and removals.
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		OnAdded:   false,
		OnRemoved: true,
		OnMajor:   true,
		OnMinor:   false,
	}
}

// GenerateAlerts produces alerts from a DiffResult based on the given config.
func GenerateAlerts(diff DiffResult, cfg AlertConfig) []Alert {
	var alerts []Alert

	if cfg.OnAdded {
		for _, dep := range diff.Added {
			alerts = append(alerts, Alert{
				Module:  dep.Module,
				Level:   AlertInfo,
				Message: fmt.Sprintf("new dependency added: %s@%s", dep.Module, dep.Version),
			})
		}
	}

	if cfg.OnRemoved {
		for _, dep := range diff.Removed {
			alerts = append(alerts, Alert{
				Module:  dep.Module,
				Level:   AlertWarning,
				Message: fmt.Sprintf("dependency removed: %s@%s", dep.Module, dep.Version),
			})
		}
	}

	for _, ch := range diff.Changed {
		major := isMajorChange(ch.OldVersion, ch.NewVersion)
		if cfg.OnMajor && major {
			alerts = append(alerts, Alert{
				Module:  ch.Module,
				Level:   AlertCritical,
				Message: fmt.Sprintf("major version change: %s %s -> %s", ch.Module, ch.OldVersion, ch.NewVersion),
			})
		} else if cfg.OnMinor && !major {
			alerts = append(alerts, Alert{
				Module:  ch.Module,
				Level:   AlertInfo,
				Message: fmt.Sprintf("version updated: %s %s -> %s", ch.Module, ch.OldVersion, ch.NewVersion),
			})
		}
	}

	return alerts
}

// FormatAlerts returns a human-readable summary of all alerts.
func FormatAlerts(alerts []Alert) string {
	if len(alerts) == 0 {
		return "no alerts"
	}
	var sb strings.Builder
	for _, a := range alerts {
		fmt.Fprintf(&sb, "[%s] %s\n", a.Level, a.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
