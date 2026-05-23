package gomod

import "fmt"

// BreakingScore represents a numeric risk score for a diff result.
type BreakingScore struct {
	Total   int
	Added   int
	Removed int
	Changed int
	Major   int
}

// String returns a human-readable summary of the score.
func (s BreakingScore) String() string {
	return fmt.Sprintf("score=%d (added=%d, removed=%d, changed=%d, major=%d)",
		s.Total, s.Added, s.Removed, s.Changed, s.Major)
}

// ScoreWeights controls how much each change type contributes to the total.
type ScoreWeights struct {
	Added   int
	Removed int
	Changed int
	Major   int
}

// DefaultWeights returns a sensible default weighting.
func DefaultWeights() ScoreWeights {
	return ScoreWeights{
		Added:   1,
		Removed: 2,
		Changed: 2,
		Major:   5,
	}
}

// ScoreDiff computes a BreakingScore for a DiffResult using the provided weights.
func ScoreDiff(diff DiffResult, weights ScoreWeights) BreakingScore {
	var s BreakingScore

	for _, c := range diff.Changes {
		switch c.Type {
		case ChangeAdded:
			s.Added++
			s.Total += weights.Added
		case ChangeRemoved:
			s.Removed++
			s.Total += weights.Removed
		case ChangeUpdated:
			if isMajorChange(c.OldVersion, c.NewVersion) {
				s.Major++
				s.Total += weights.Major
			} else {
				s.Changed++
				s.Total += weights.Changed
			}
		}
	}

	return s
}

// RiskLevel returns a qualitative label based on the total score.
func RiskLevel(score BreakingScore) string {
	switch {
	case score.Total == 0:
		return "none"
	case score.Total <= 5:
		return "low"
	case score.Total <= 15:
		return "medium"
	default:
		return "high"
	}
}
