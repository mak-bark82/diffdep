package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveDeprecationReport writes a DeprecationReport to disk as JSON.
func SaveDeprecationReport(dir string, r DeprecationReport) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("deprecate store: mkdir: %w", err)
	}
	path := filepath.Join(dir, "deprecation_report.json")
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("deprecate store: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("deprecate store: write: %w", err)
	}
	return nil
}

// LoadDeprecationReport reads a DeprecationReport from disk.
func LoadDeprecationReport(dir string) (DeprecationReport, error) {
	path := filepath.Join(dir, "deprecation_report.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DeprecationReport{}, fmt.Errorf("deprecate store: file not found: %s", path)
		}
		return DeprecationReport{}, fmt.Errorf("deprecate store: read: %w", err)
	}
	var r DeprecationReport
	if err := json.Unmarshal(data, &r); err != nil {
		return DeprecationReport{}, fmt.Errorf("deprecate store: unmarshal: %w", err)
	}
	return r, nil
}
