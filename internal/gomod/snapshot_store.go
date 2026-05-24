package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const snapshotFileExt = ".snapshot.json"

// SaveSnapshot persists a Snapshot to dir/<branch>.snapshot.json.
func SaveSnapshot(dir string, s Snapshot) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir %s: %w", dir, err)
	}
	path := filepath.Join(dir, sanitizeBranch(s.Branch)+snapshotFileExt)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from dir/<branch>.snapshot.json.
func LoadSnapshot(dir, branch string) (Snapshot, error) {
	path := filepath.Join(dir, sanitizeBranch(branch)+snapshotFileExt)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, fmt.Errorf("snapshot: no snapshot for branch %q", branch)
		}
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode %s: %w", path, err)
	}
	return s, nil
}

// sanitizeBranch replaces path separators so branch names are safe as filenames.
func sanitizeBranch(branch string) string {
	r := []rune(branch)
	for i, c := range r {
		if c == '/' || c == '\\' {
			r[i] = '_'
		}
	}
	return string(r)
}
