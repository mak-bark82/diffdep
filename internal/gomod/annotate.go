package gomod

import (
	"fmt"
	"strings"
)

// AnnotationKind classifies the nature of a dependency change.
type AnnotationKind string

const (
	AnnotationBreaking  AnnotationKind = "breaking"
	AnnotationUpgrade   AnnotationKind = "upgrade"
	AnnotationDowngrade AnnotationKind = "downgrade"
	AnnotationNew       AnnotationKind = "new"
	AnnotationRemoved   AnnotationKind = "removed"
)

// Annotation attaches a human-readable note and classification to a diff entry.
type Annotation struct {
	Module string
	Kind   AnnotationKind
	Note   string
}

// AnnotateDiff inspects a DiffResult and returns an Annotation for each changed dependency.
func AnnotateDiff(diff DiffResult) []Annotation {
	var annotations []Annotation

	for _, e := range diff.Added {
		annotations = append(annotations, Annotation{
			Module: e.Module,
			Kind:   AnnotationNew,
			Note:   fmt.Sprintf("newly added at %s", e.Version),
		})
	}

	for _, e := range diff.Removed {
		annotations = append(annotations, Annotation{
			Module: e.Module,
			Kind:   AnnotationRemoved,
			Note:   fmt.Sprintf("removed (was %s)", e.Version),
		})
	}

	for _, e := range diff.Changed {
		kind, note := classifyChange(e)
		annotations = append(annotations, Annotation{
			Module: e.Module,
			Kind:   kind,
			Note:   note,
		})
	}

	return annotations
}

func classifyChange(e DiffEntry) (AnnotationKind, string) {
	if isMajorChange(e.OldVersion, e.NewVersion) {
		return AnnotationBreaking, fmt.Sprintf("major version bump %s → %s", e.OldVersion, e.NewVersion)
	}
	if isDowngrade(e.OldVersion, e.NewVersion) {
		return AnnotationDowngrade, fmt.Sprintf("downgraded %s → %s", e.OldVersion, e.NewVersion)
	}
	return AnnotationUpgrade, fmt.Sprintf("upgraded %s → %s", e.OldVersion, e.NewVersion)
}

// FormatAnnotations renders annotations as a plain-text list.
func FormatAnnotations(annotations []Annotation) string {
	if len(annotations) == 0 {
		return "no annotations\n"
	}
	var sb strings.Builder
	for _, a := range annotations {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", a.Kind, a.Module, a.Note))
	}
	return sb.String()
}
