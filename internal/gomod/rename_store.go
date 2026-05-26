package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveRenameReport persists a RenameReport to a JSON file at the given path.
func SaveRenameReport(path string, r RenameReport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("rename store: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("rename store: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("rename store: encode: %w", err)
	}
	return nil
}

// LoadRenameReport reads a RenameReport from a JSON file at the given path.
func LoadRenameReport(path string) (RenameReport, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return RenameReport{}, fmt.Errorf("rename store: file not found: %s", path)
		}
		return RenameReport{}, fmt.Errorf("rename store: open: %w", err)
	}
	defer f.Close()
	var r RenameReport
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return RenameReport{}, fmt.Errorf("rename store: decode: %w", err)
	}
	return r, nil
}
