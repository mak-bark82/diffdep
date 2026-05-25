package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveLicenseReport persists a LicenseReport to a JSON file under dir.
func SaveLicenseReport(dir string, report *LicenseReport) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("license store: mkdir: %w", err)
	}
	path := filepath.Join(dir, "license_report.json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("license store: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("license store: encode: %w", err)
	}
	return nil
}

// LoadLicenseReport reads a previously saved LicenseReport from dir.
func LoadLicenseReport(dir string) (*LicenseReport, error) {
	path := filepath.Join(dir, "license_report.json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("license store: no report found at %s", path)
		}
		return nil, fmt.Errorf("license store: open: %w", err)
	}
	defer f.Close()
	var report LicenseReport
	if err := json.NewDecoder(f).Decode(&report); err != nil {
		return nil, fmt.Errorf("license store: decode: %w", err)
	}
	return &report, nil
}
