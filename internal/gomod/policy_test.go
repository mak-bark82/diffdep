package gomod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func samplePolicyDiff() []DiffEntry {
	return []DiffEntry{
		{Module: "github.com/foo/bar", ChangeType: ChangeAdded, NewVersion: "v1.0.0"},
		{Module: "github.com/foo/baz", ChangeType: ChangeAdded, NewVersion: "v2.0.0"},
		{Module: "github.com/old/dep", ChangeType: ChangeRemoved, OldVersion: "v1.0.0"},
		{Module: "github.com/upgrade/me", ChangeType: ChangeUpdated, OldVersion: "v1.3.0", NewVersion: "v2.0.0"},
	}
}

func TestEnforcePolicy_NoViolations(t *testing.T) {
	p := DefaultPolicy()
	violations := EnforcePolicy(p, samplePolicyDiff())
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestEnforcePolicy_MaxAdded(t *testing.T) {
	p := DefaultPolicy()
	p.MaxAdded = 1
	violations := EnforcePolicy(p, samplePolicyDiff())
	if len(violations) != 1 || violations[0].Rule != "max_added" {
		t.Fatalf("expected max_added violation, got %v", violations)
	}
}

func TestEnforcePolicy_MaxMajor(t *testing.T) {
	p := DefaultPolicy()
	p.MaxMajor = 0
	violations := EnforcePolicy(p, samplePolicyDiff())
	if len(violations) != 1 || violations[0].Rule != "max_major" {
		t.Fatalf("expected max_major violation, got %v", violations)
	}
}

func TestEnforcePolicy_DenyList(t *testing.T) {
	p := DefaultPolicy()
	p.DenyList = []string{"github.com/foo"}
	violations := EnforcePolicy(p, samplePolicyDiff())
	if len(violations) != 2 {
		t.Fatalf("expected 2 deny_list violations, got %d", len(violations))
	}
	for _, v := range violations {
		if v.Rule != "deny_list" {
			t.Errorf("unexpected rule %q", v.Rule)
		}
	}
}

func TestLoadPolicy_Missing(t *testing.T) {
	p, err := LoadPolicy("/nonexistent/policy.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if p.MaxAdded != -1 {
		t.Errorf("expected default policy")
	}
}

func TestLoadPolicy_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	data, _ := json.Marshal(Policy{MaxAdded: 5, MaxMajor: 1, DenyList: []string{"bad/pkg"}})
	os.WriteFile(path, data, 0644)

	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxAdded != 5 || p.MaxMajor != 1 || len(p.DenyList) != 1 {
		t.Errorf("policy not loaded correctly: %+v", p)
	}
}

func TestLoadPolicy_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
