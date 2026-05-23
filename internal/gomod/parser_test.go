package gomod

import (
	"strings"
	"testing"
)

func TestParseGoMod(t *testing.T) {
	input := `module github.com/example/app

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/sync v0.6.0 // indirect
)

require github.com/stretchr/testify v1.8.4
`
	deps, err := ParseGoMod(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %d", len(deps))
	}

	cases := []struct {
		module   string
		version  string
		indirect bool
	}{
		{"github.com/pkg/errors", "v0.9.1", false},
		{"golang.org/x/sync", "v0.6.0", true},
		{"github.com/stretchr/testify", "v1.8.4", false},
	}
	for i, tc := range cases {
		if deps[i].Module != tc.module {
			t.Errorf("dep[%d].Module = %q, want %q", i, deps[i].Module, tc.module)
		}
		if deps[i].Version != tc.version {
			t.Errorf("dep[%d].Version = %q, want %q", i, deps[i].Version, tc.version)
		}
		if deps[i].Indirect != tc.indirect {
			t.Errorf("dep[%d].Indirect = %v, want %v", i, deps[i].Indirect, tc.indirect)
		}
	}
}

func TestParseGoSum(t *testing.T) {
	input := `github.com/pkg/errors v0.9.1 h1:abc123==
github.com/pkg/errors v0.9.1/go.mod h1:def456==
golang.org/x/sync v0.6.0 h1:xyz789==
`
	deps, err := ParseGoSum(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 2 {
		t.Fatalf("expected 2 unique modules, got %d", len(deps))
	}
	if v := deps["github.com/pkg/errors"]; v != "v0.9.1" {
		t.Errorf("pkg/errors version = %q, want v0.9.1", v)
	}
	if v := deps["golang.org/x/sync"]; v != "v0.6.0" {
		t.Errorf("x/sync version = %q, want v0.6.0", v)
	}
}

func TestParseGoMod_Empty(t *testing.T) {
	deps, err := ParseGoMod(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Fatalf("expected 0 deps, got %d", len(deps))
	}
}
