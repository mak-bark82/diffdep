package gomod

// Dependency represents a single module dependency with its version.
type Dependency struct {
	Module  string
	Version string
}

// DependencyChange represents a version change for a single module between
// two branches or snapshots.
type DependencyChange struct {
	Module     string
	OldVersion string
	NewVersion string
}

// DiffResult holds the categorised differences between two dependency sets.
type DiffResult struct {
	Added   []Dependency
	Removed []Dependency
	Changed []DependencyChange
}

// IsEmpty returns true when the DiffResult contains no changes.
func (d DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}
