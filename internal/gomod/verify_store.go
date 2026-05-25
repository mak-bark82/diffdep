package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveVerifyReport persists a VerifyReport to disk as JSON.
func SaveVerifyReport(dir string, report VerifyReport) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("verify store: mkdir: %w", err)
	}
	path := filepath.Join(dir, "verify_report.json")
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("verify store: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("verify store: write: %w", err)
	}
	return nil
}

// LoadVerifyReport reads a VerifyReport from disk.
func LoadVerifyReport(dir string) (VerifyReport, error) {
	path := filepath.Join(dir, "verify_report.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return VerifyReport{}, fmt.Errorf("verify store: report not found at %s", path)
		}
		return VerifyReport{}, fmt.Errorf("verify store: read: %w", err)
	}
	var report VerifyReport
	if err := json.Unmarshal(data, &report); err != nil {
		return VerifyReport{}, fmt.Errorf("verify store: unmarshal: %w", err)
	}
	return report, nil
}
