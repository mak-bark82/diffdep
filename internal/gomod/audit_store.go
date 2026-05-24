package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveAuditReport persists an AuditReport to a JSON file.
func SaveAuditReport(dir string, report AuditReport) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("audit: create dir: %w", err)
	}
	path := filepath.Join(dir, "audit.json")
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}

// LoadAuditReport reads an AuditReport from a JSON file.
func LoadAuditReport(dir string) (AuditReport, error) {
	path := filepath.Join(dir, "audit.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AuditReport{}, fmt.Errorf("audit: file not found: %s", path)
		}
		return AuditReport{}, fmt.Errorf("audit: read: %w", err)
	}
	var report AuditReport
	if err := json.Unmarshal(data, &report); err != nil {
		return AuditReport{}, fmt.Errorf("audit: unmarshal: %w", err)
	}
	return report, nil
}
