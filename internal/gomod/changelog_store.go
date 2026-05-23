package gomod

import (
	"encoding/json"
	"errors"
	"os"
)

const defaultChangelogPath = ".diffdep-changelog.json"

// SaveChangelog persists a Changelog to disk as JSON.
func SaveChangelog(path string, c *Changelog) error {
	if path == "" {
		path = defaultChangelogPath
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal changelog: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadChangelog reads a Changelog from disk. Returns an empty Changelog if the
// file does not exist yet.
func LoadChangelog(path string) (*Changelog, error) {
	if path == "" {
		path = defaultChangelogPath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Changelog{}, nil
		}
		return nil, fmt.Errorf("read changelog: %w", err)
	}
	var c Changelog
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal changelog: %w", err)
	}
	return &c, nil
}
