package gomod

// Dependency represents a single module dependency with its version.
type Dependency struct {
	Module  string
	Version string
}

// DiffEntry represents a change in a dependency between two branches.
type DiffEntry struct {
	Module  string
	OldVersion string
	NewVersion string
	ChangeType ChangeType
}

// ChangeType indicates how a dependency changed.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// ReportFormat specifies the output format for reports.
type ReportFormat string

const (
	TextFormat     ReportFormat = "text"
	MarkdownFormat ReportFormat = "markdown"
)
