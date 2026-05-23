package gomod

import "fmt"

// ChangeKind describes the type of version change between two branches.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Updated ChangeKind = "updated"
)

// Change represents a single dependency change between two branches.
type Change struct {
	Module  string
	Kind    ChangeKind
	OldVer  string
	NewVer  string
}

// String returns a human-readable representation of the change.
func (c Change) String() string {
	switch c.Kind {
	case Added:
		return fmt.Sprintf("+ %s %s", c.Module, c.NewVer)
	case Removed:
		return fmt.Sprintf("- %s %s", c.Module, c.OldVer)
	case Updated:
		return fmt.Sprintf("~ %s %s -> %s", c.Module, c.OldVer, c.NewVer)
	}
	return ""
}

// DiffDependencies compares two sets of dependencies (base vs head) and
// returns a list of changes. Each map is keyed by module path with version as value.
func DiffDependencies(base, head map[string]string) []Change {
	var changes []Change

	// Detect removed and updated modules.
	for mod, baseVer := range base {
		headVer, exists := head[mod]
		if !exists {
			changes = append(changes, Change{
				Module: mod,
				Kind:   Removed,
				OldVer: baseVer,
			})
		} else if baseVer != headVer {
			changes = append(changes, Change{
				Module: mod,
				Kind:   Updated,
				OldVer: baseVer,
				NewVer: headVer,
			})
		}
	}

	// Detect added modules.
	for mod, headVer := range head {
		if _, exists := base[mod]; !exists {
			changes = append(changes, Change{
				Module: mod,
				Kind:   Added,
				NewVer: headVer,
			})
		}
	}

	return changes
}

// DepsToMap converts a slice of Dependency into a module->version map,
// suitable for use with DiffDependencies.
func DepsToMap(deps []Dependency) map[string]string {
	m := make(map[string]string, len(deps))
	for _, d := range deps {
		m[d.Module] = d.Version
	}
	return m
}
