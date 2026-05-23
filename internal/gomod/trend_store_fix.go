package gomod

import "fmt"

// ensure fmt is used by trend_store.go via this shim.
// This file exists because trend_store.go references fmt.Errorf
// and Go requires the import in the same file.
// Re-export a helper used internally.
func trendFmtErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
