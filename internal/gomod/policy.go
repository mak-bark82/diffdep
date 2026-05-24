package gomod

import (
	"encoding/json"
	"fmt"
	"os"
)

// Policy defines rules that determine whether a diff should cause a failure.
type Policy struct {
	MaxAdded      int      `json:"max_added"`
	MaxRemoved    int      `json:"max_removed"`
	MaxMajor      int      `json:"max_major"`
	DenyList      []string `json:"deny_list"`
	RequireMajor  bool     `json:"require_major_review"`
}

// PolicyViolation describes a single rule that was broken.
type PolicyViolation struct {
	Rule    string
	Message string
}

func (v PolicyViolation) String() string {
	return fmt.Sprintf("[%s] %s", v.Rule, v.Message)
}

// DefaultPolicy returns a permissive default policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAdded:   -1,
		MaxRemoved: -1,
		MaxMajor:   -1,
	}
}

// LoadPolicy reads a policy from a JSON file.
func LoadPolicy(path string) (Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPolicy(), nil
		}
		return Policy{}, fmt.Errorf("read policy: %w", err)
	}
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return Policy{}, fmt.Errorf("parse policy: %w", err)
	}
	return p, nil
}

// EnforcePolicy checks the diff against the policy and returns any violations.
func EnforcePolicy(p Policy, diff []DiffEntry) []PolicyViolation {
	var violations []PolicyViolation

	added, removed, major := 0, 0, 0
	for _, d := range diff {
		switch d.ChangeType {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeUpdated:
			if isMajorChange(d.OldVersion, d.NewVersion) {
				major++
			}
		}
		for _, denied := range p.DenyList {
			if matchesPrefix(d.Module, denied) {
				violations = append(violations, PolicyViolation{
					Rule:    "deny_list",
					Message: fmt.Sprintf("module %q is on the deny list", d.Module),
				})
			}
		}
	}

	if p.MaxAdded >= 0 && added > p.MaxAdded {
		violations = append(violations, PolicyViolation{
			Rule:    "max_added",
			Message: fmt.Sprintf("added %d dependencies, limit is %d", added, p.MaxAdded),
		})
	}
	if p.MaxRemoved >= 0 && removed > p.MaxRemoved {
		violations = append(violations, PolicyViolation{
			Rule:    "max_removed",
			Message: fmt.Sprintf("removed %d dependencies, limit is %d", removed, p.MaxRemoved),
		})
	}
	if p.MaxMajor >= 0 && major > p.MaxMajor {
		violations = append(violations, PolicyViolation{
			Rule:    "max_major",
			Message: fmt.Sprintf("%d major version bumps, limit is %d", major, p.MaxMajor),
		})
	}
	if p.RequireMajor && major == 0 {
		violations = append(violations, PolicyViolation{
			Rule:    "require_major_review",
			Message: "policy requires at least one major change to trigger review, but none found",
		})
	}

	return violations
}
