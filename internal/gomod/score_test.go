package gomod

import (
	"testing"
)

func TestScoreDiff_Empty(t *testing.T) {
	diff := DiffResult{}
	score := ScoreDiff(diff, DefaultWeights())
	if score.Total != 0 {
		t.Errorf("expected total 0, got %d", score.Total)
	}
	if RiskLevel(score) != "none" {
		t.Errorf("expected risk 'none', got %s", RiskLevel(score))
	}
}

func TestScoreDiff_AddedOnly(t *testing.T) {
	diff := DiffResult{
		Changes: []DependencyChange{
			{Module: "github.com/foo/bar", Type: ChangeAdded, NewVersion: "v1.0.0"},
			{Module: "github.com/foo/baz", Type: ChangeAdded, NewVersion: "v2.0.0"},
		},
	}
	weights := DefaultWeights()
	score := ScoreDiff(diff, weights)
	if score.Added != 2 {
		t.Errorf("expected 2 added, got %d", score.Added)
	}
	if score.Total != 2*weights.Added {
		t.Errorf("expected total %d, got %d", 2*weights.Added, score.Total)
	}
}

func TestScoreDiff_MajorChange(t *testing.T) {
	diff := DiffResult{
		Changes: []DependencyChange{
			{Module: "github.com/foo/bar", Type: ChangeUpdated, OldVersion: "v1.2.3", NewVersion: "v2.0.0"},
		},
	}
	weights := DefaultWeights()
	score := ScoreDiff(diff, weights)
	if score.Major != 1 {
		t.Errorf("expected 1 major change, got %d", score.Major)
	}
	if score.Total != weights.Major {
		t.Errorf("expected total %d, got %d", weights.Major, score.Total)
	}
	if RiskLevel(score) != "low" {
		t.Errorf("expected risk 'low', got %s", RiskLevel(score))
	}
}

func TestScoreDiff_HighRisk(t *testing.T) {
	changes := make([]DependencyChange, 4)
	for i := range changes {
		changes[i] = DependencyChange{
			Module:     "github.com/foo/pkg",
			Type:       ChangeUpdated,
			OldVersion: "v1.0.0",
			NewVersion: "v2.0.0",
		}
	}
	diff := DiffResult{Changes: changes}
	score := ScoreDiff(diff, DefaultWeights())
	if RiskLevel(score) != "high" {
		t.Errorf("expected risk 'high', got %s", RiskLevel(score))
	}
}

func TestBreakingScore_String(t *testing.T) {
	s := BreakingScore{Total: 10, Added: 1, Removed: 2, Changed: 1, Major: 1}
	str := s.String()
	if str == "" {
		t.Error("expected non-empty string from BreakingScore.String()")
	}
}
