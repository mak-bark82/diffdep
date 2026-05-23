package gomod

// Dependency represents a single module dependency with its resolved version.
type Dependency struct {
	Path    string
	Version string
}

// DiffResult holds the categorised changes between two dependency sets.
type DiffResult struct {
	Added   []Dependency
	Removed []Dependency
	Changed []DependencyChange
}

// DependencyChange records a version transition for a single module.
type DependencyChange struct {
	Path       string
	OldVersion string
	NewVersion string
}

// HasChanges reports whether the diff contains any additions, removals, or
// version changes.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
