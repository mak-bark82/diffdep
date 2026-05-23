package gomod

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Dependency represents a single module dependency with its version.
type Dependency struct {
	Module  string
	Version string
	Indirect bool
}

// ParseGoSum parses a go.sum file content and returns a map of module@version entries.
func ParseGoSum(r io.Reader) (map[string]string, error) {
	deps := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		modVer := parts[0]
		segments := strings.SplitN(modVer, "@", 2)
		if len(segments) != 2 {
			continue
		}
		mod, ver := segments[0], segments[1]
		// go.sum may list /go.mod entries; strip that suffix from version
		ver = strings.TrimSuffix(ver, "/go.mod")
		deps[mod] = ver
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning go.sum: %w", err)
	}
	return deps, nil
}

// ParseGoMod parses a go.mod file content and returns a slice of Dependency.
func ParseGoMod(r io.Reader) ([]Dependency, error) {
	var deps []Dependency
	scanner := bufio.NewScanner(r)
	inRequire := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if strings.HasPrefix(line, "require ") {
			line = strings.TrimPrefix(line, "require ")
			dep, err := parseDependencyLine(line)
			if err == nil {
				deps = append(deps, dep)
			}
			continue
		}
		if inRequire {
			dep, err := parseDependencyLine(line)
			if err == nil {
				deps = append(deps, dep)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning go.mod: %w", err)
	}
	return deps, nil
}

func parseDependencyLine(line string) (Dependency, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Dependency{}, fmt.Errorf("invalid dependency line: %q", line)
	}
	dep := Dependency{
		Module:   parts[0],
		Version:  parts[1],
		Indirect: len(parts) >= 3 && strings.Contains(line, "// indirect"),
	}
	return dep, nil
}
