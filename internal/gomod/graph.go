package gomod

import (
	"fmt"
	"sort"
	"strings"
)

// GraphNode represents a dependency node in the graph.
type GraphNode struct {
	Module  string
	Version string
	Edges   []string // modules this node depends on
}

// DependencyGraph holds a directed graph of module dependencies.
type DependencyGraph struct {
	Nodes map[string]*GraphNode
}

// NewDependencyGraph builds a DependencyGraph from a slice of Dependencies.
func NewDependencyGraph(deps []Dependency) *DependencyGraph {
	g := &DependencyGraph{
		Nodes: make(map[string]*GraphNode),
	}
	for _, d := range deps {
		g.Nodes[d.Module] = &GraphNode{
			Module:  d.Module,
			Version: d.Version,
			Edges:   []string{},
		}
	}
	return g
}

// AddEdge records that module `from` depends on module `to`.
func (g *DependencyGraph) AddEdge(from, to string) error {
	node, ok := g.Nodes[from]
	if !ok {
		return fmt.Errorf("module %q not found in graph", from)
	}
	node.Edges = append(node.Edges, to)
	return nil
}

// Affected returns all modules that transitively depend on the given module.
func (g *DependencyGraph) Affected(module string) []string {
	visited := map[string]bool{}
	var result []string
	for name, node := range g.Nodes {
		if name == module {
			continue
		}
		if g.reachable(name, module, visited) {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func (g *DependencyGraph) reachable(from, target string, cache map[string]bool) bool {
	key := from + "->" + target
	if v, ok := cache[key]; ok {
		return v
	}
	node, ok := g.Nodes[from]
	if !ok {
		cache[key] = false
		return false
	}
	for _, edge := range node.Edges {
		if edge == target || g.reachable(edge, target, cache) {
			cache[key] = true
			return true
		}
	}
	cache[key] = false
	return false
}

// FormatGraph returns a simple text representation of the graph.
func FormatGraph(g *DependencyGraph) string {
	var sb strings.Builder
	keys := make([]string, 0, len(g.Nodes))
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		node := g.Nodes[k]
		if len(node.Edges) == 0 {
			fmt.Fprintf(&sb, "%s@%s\n", node.Module, node.Version)
		} else {
			sort.Strings(node.Edges)
			fmt.Fprintf(&sb, "%s@%s -> [%s]\n", node.Module, node.Version, strings.Join(node.Edges, ", "))
		}
	}
	return sb.String()
}
