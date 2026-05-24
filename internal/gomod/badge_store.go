package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveBadge persists BadgeData as JSON to the given directory.
func SaveBadge(dir string, b BadgeData) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("badge: mkdir %s: %w", dir, err)
	}
	path := filepath.Join(dir, "badge.json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("badge: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("badge: encode: %w", err)
	}
	return nil
}

// LoadBadge reads BadgeData from a JSON file in the given directory.
func LoadBadge(dir string) (BadgeData, error) {
	path := filepath.Join(dir, "badge.json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return BadgeData{}, fmt.Errorf("badge: file not found: %s", path)
		}
		return BadgeData{}, fmt.Errorf("badge: open %s: %w", path, err)
	}
	defer f.Close()
	var b BadgeData
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return BadgeData{}, fmt.Errorf("badge: decode: %w", err)
	}
	return b, nil
}
