package gomod

import (
	"encoding/json"
	"errors"
	"os"
)

// SaveTrend persists a Trend to the given file path as JSON.
func SaveTrend(path string, trend *Trend) error {
	if trend == nil {
		return errors.New("trend must not be nil")
	}
	data, err := json.MarshalIndent(trend, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal trend: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadTrend reads a Trend from the given file path.
// If the file does not exist, an empty Trend is returned.
func LoadTrend(path string) (*Trend, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Trend{}, nil
		}
		return nil, fmt.Errorf("read trend file: %w", err)
	}
	var trend Trend
	if err := json.Unmarshal(data, &trend); err != nil {
		return nil, fmt.Errorf("unmarshal trend: %w", err)
	}
	return &trend, nil
}
