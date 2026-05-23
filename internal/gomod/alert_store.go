package gomod

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveAlerts writes a slice of alerts to a JSON file at the given path.
func SaveAlerts(path string, alerts []Alert) error {
	data, err := json.MarshalIndent(alerts, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal alerts: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write alerts file: %w", err)
	}
	return nil
}

// LoadAlerts reads a slice of alerts from a JSON file at the given path.
// Returns an empty slice if the file does not exist.
func LoadAlerts(path string) ([]Alert, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Alert{}, nil
		}
		return nil, fmt.Errorf("read alerts file: %w", err)
	}
	var alerts []Alert
	if err := json.Unmarshal(data, &alerts); err != nil {
		return nil, fmt.Errorf("unmarshal alerts: %w", err)
	}
	return alerts, nil
}
