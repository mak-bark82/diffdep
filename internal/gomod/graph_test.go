package gomod

import (
	"strings"
	"testing"
)

func sampleDepsForGraph() []Dependency {
	return []Dependency{
		{Module: "github.com/a/core", Version: "v1.0.0"},
		{Module: "github.com/b/util", Version: "v2.1.0"},
		{Module: "github.com/c/app", Version: "v0.5.0"},
	}
}

func TestNewDependencyGraph_NodeCount(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	if len(g.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.Nodes))
	}
}

func TestAddEdge_Valid(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	if err := g.AddEdge("github.com/c/app", "github.com/a/core"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := g.Nodes["github.com/c/app"]
	if len(node.Edges) != 1 || node.Edges[0] != "github.com/a/core" {
		t.Errorf("edge not recorded correctly: %v", node.Edges)
	}
}

func TestAddEdge_UnknownModule(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	err := g.AddEdge("github.com/unknown/pkg", "github.com/a/core")
	if err == nil {
		t.Fatal("expected error for unknown module, got nil")
	}
}

func TestAffected_DirectDependents(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	_ = g.AddEdge("github.com/c/app", "github.com/a/core")
	_ = g.AddEdge("github.com/b/util", "github.com/a/core")

	affected := g.Affected("github.com/a/core")
	if len(affected) != 2 {
		t.Fatalf("expected 2 affected modules, got %d: %v", len(affected), affected)
	}
}

func TestAffected_NoDependent(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	affected := g.Affected("github.com/a/core")
	if len(affected) != 0 {
		t.Errorf("expected no affected modules, got %v", affected)
	}
}

func TestAffected_Transitive(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	_ = g.AddEdge("github.com/b/util", "github.com/a/core")
	_ = g.AddEdge("github.com/c/app", "github.com/b/util")

	affected := g.Affected("github.com/a/core")
	if len(affected) != 2 {
		t.Fatalf("expected 2 transitively affected modules, got %d: %v", len(affected), affected)
	}
}

func TestFormatGraph_ContainsModules(t *testing.T) {
	g := NewDependencyGraph(sampleDepsForGraph())
	_ = g.AddEdge("github.com/c/app", "github.com/a/core")
	out := FormatGraph(g)
	if !strings.Contains(out, "github.com/c/app") {
		t.Errorf("expected output to contain github.com/c/app, got:\n%s", out)
	}
	if !strings.Contains(out, "github.com/a/core") {
		t.Errorf("expected output to contain github.com/a/core, got:\n%s", out)
	}
}
