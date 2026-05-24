package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const lockFileName = "diffdep.lock.json"

// SaveLockFile writes a LockFile to the given directory.
func SaveLockFile(dir string, lf *LockFile) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("lock: mkdir %s: %w", dir, err)
	}
	path := filepath.Join(dir, lockFileName)
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("lock: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("lock: write %s: %w", path, err)
	}
	return nil
}

// LoadLockFile reads a LockFile from the given directory.
// Returns nil, nil if the file does not exist.
func LoadLockFile(dir string) (*LockFile, error) {
	path := filepath.Join(dir, lockFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("lock: read %s: %w", path, err)
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("lock: parse %s: %w", path, err)
	}
	return &lf, nil
}
