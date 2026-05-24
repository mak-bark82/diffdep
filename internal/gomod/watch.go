package gomod

import (
	"fmt"
	"strings"
	"time"
)

// WatchEvent represents a detected change during a watch cycle.
type WatchEvent struct {
	Branch    string
	Timestamp time.Time
	Diff      []DiffEntry
	Alerts    []Alert
}

// WatchConfig configures polling behavior for the watch command.
type WatchConfig struct {
	Branch       string
	IntervalSecs int
	AlertConfig  AlertConfig
	OnChange     func(WatchEvent)
}

// DefaultWatchConfig returns a WatchConfig with sensible defaults.
func DefaultWatchConfig(branch string) WatchConfig {
	return WatchConfig{
		Branch:       branch,
		IntervalSecs: 60,
		AlertConfig:  DefaultAlertConfig(),
	}
}

// RunWatch polls dependency state at a fixed interval and fires OnChange
// whenever the diff against the provided baseline changes.
func RunWatch(cfg WatchConfig, loader func(branch string) ([]Dependency, error), stop <-chan struct{}) error {
	if cfg.IntervalSecs <= 0 {
		return fmt.Errorf("watch: interval must be positive")
	}
	if loader == nil {
		return fmt.Errorf("watch: loader must not be nil")
	}

	ticker := time.NewTicker(time.Duration(cfg.IntervalSecs) * time.Second)
	defer ticker.Stop()

	var lastDiffKey string

	check := func() {
		deps, err := loader(cfg.Branch)
		if err != nil {
			return
		}
		snap := NewSnapshot(cfg.Branch, deps)
		prev, err := LoadSnapshot(cfg.Branch)
		if err != nil {
			_ = SaveSnapshot(snap)
			return
		}
		diff := DiffSnapshot(prev, snap)
		key := diffKey(diff)
		if key == lastDiffKey {
			return
		}
		lastDiffKey = key
		_ = SaveSnapshot(snap)
		if cfg.OnChange != nil {
			alerts := GenerateAlerts(diff, cfg.AlertConfig)
			cfg.OnChange(WatchEvent{
				Branch:    cfg.Branch,
				Timestamp: time.Now(),
				Diff:      diff,
				Alerts:    alerts,
			})
		}
	}

	check()
	for {
		select {
		case <-ticker.C:
			check()
		case <-stop:
			return nil
		}
	}
}

// diffKey builds a lightweight string key from a diff slice for change detection.
func diffKey(diff []DiffEntry) string {
	parts := make([]string, 0, len(diff))
	for _, d := range diff {
		parts = append(parts, fmt.Sprintf("%s@%s->%s", d.Module, d.OldVersion, d.NewVersion))
	}
	return strings.Join(parts, "|")
}
