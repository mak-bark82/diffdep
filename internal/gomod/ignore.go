package gomod

import (
	"bufio"
	"os"
	"strings"
)

// IgnoreList holds a set of module path prefixes to exclude from diffs.
type IgnoreList struct {
	prefixes []string
}

// LoadIgnoreFile reads a .diffdepignore file where each non-blank,
// non-comment line is treated as a module path prefix to ignore.
func LoadIgnoreFile(path string) (*IgnoreList, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &IgnoreList{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var prefixes []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		prefixes = append(prefixes, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &IgnoreList{prefixes: prefixes}, nil
}

// NewIgnoreList creates an IgnoreList from a slice of prefix strings.
func NewIgnoreList(prefixes []string) *IgnoreList {
	cp := make([]string, len(prefixes))
	copy(cp, prefixes)
	return &IgnoreList{prefixes: cp}
}

// ShouldIgnore returns true if the given module path matches any prefix
// in the ignore list.
func (il *IgnoreList) ShouldIgnore(module string) bool {
	for _, p := range il.prefixes {
		if strings.HasPrefix(module, p) {
			return true
		}
	}
	return false
}

// ApplyIgnore removes entries from a DiffResult whose module path is
// covered by the ignore list.
func ApplyIgnore(diff DiffResult, il *IgnoreList) DiffResult {
	if il == nil || len(il.prefixes) == 0 {
		return diff
	}
	filtered := DiffResult{}
	for _, e := range diff {
		if !il.ShouldIgnore(e.Module) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
