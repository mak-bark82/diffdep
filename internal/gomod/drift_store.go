package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveDriftReport persists a DriftReport to the given directory as JSON.
func SaveDriftReport(dir string, report DriftReport) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("drift store: mkdir: %w", err)
	}
	path := filepath.Join(dir, "drift.json")
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("drift store: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("drift store: write: %w", err)
	}
	return nil
}

// LoadDriftReport reads a previously saved DriftReport from dir.
func LoadDriftReport(dir string) (DriftReport, error) {
	path := filepath.Join(dir, "drift.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DriftReport{}, fmt.Errorf("drift store: no report found at %s", path)
		}
		return DriftReport{}, fmt.Errorf("drift store: read: %w", err)
	}
	var report DriftReport
	if err := json.Unmarshal(data, &report); err != nil {
		return DriftReport{}, fmt.Errorf("drift store: unmarshal: %w", err)
	}
	return report, nil
}
