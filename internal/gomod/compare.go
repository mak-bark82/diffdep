package gomod

import "fmt"

// BaselineComparison holds the result of comparing current deps against a baseline.
type BaselineComparison struct {
	Baseline *Baseline
	Diff     []DiffEntry
	Summary  Summary
}

// CompareAgainstBaseline diffs current dependencies against a saved baseline.
func CompareAgainstBaseline(baselinePath string, current []Dependency) (*BaselineComparison, error) {
	b, err := LoadBaseline(baselinePath)
	if err != nil {
		return nil, fmt.Errorf("load baseline: %w", err)
	}

	diff := DiffDependencies(b.Deps, current)
	sum := Summarize(diff)

	return &BaselineComparison{
		Baseline: b,
		Diff:     diff,
		Summary:  sum,
	}, nil
}

// HasBreakingChanges returns true if the comparison contains any major version changes.
func (c *BaselineComparison) HasBreakingChanges() bool {
	for _, e := range c.Diff {
		if e.ChangeType == Changed && isMajorChange(e.OldVersion, e.NewVersion) {
			return true
		}
	}
	return false
}
