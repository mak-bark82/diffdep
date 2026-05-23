package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of dependencies at a point in time.
type Baseline struct {
	Branch    string            `json:"branch"`
	CreatedAt time.Time         `json:"created_at"`
	Deps      []Dependency      `json:"deps"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// SaveBaseline writes a dependency snapshot to a JSON file.
func SaveBaseline(path string, branch string, deps []Dependency) error {
	b := Baseline{
		Branch:    branch,
		CreatedAt: time.Now().UTC(),
		Deps:      deps,
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create baseline file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// LoadBaseline reads a previously saved dependency snapshot from a JSON file.
func LoadBaseline(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open baseline file: %w", err)
	}
	defer f.Close()

	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode baseline: %w", err)
	}
	return &b, nil
}
