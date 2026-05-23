package gomod_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/diffdep/internal/gomod"
)

func TestSaveAndLoadBaseline(t *testing.T) {
	deps := []gomod.Dependency{
		{Path: "github.com/foo/bar", Version: "v1.2.3"},
		{Path: "github.com/baz/qux", Version: "v0.4.1"},
	}

	tmp := t.TempDir()
	path := filepath.Join(tmp, "baseline.json")

	if err := gomod.SaveBaseline(path, "main", deps); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	b, err := gomod.LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if b.Branch != "main" {
		t.Errorf("expected branch main, got %s", b.Branch)
	}
	if len(b.Deps) != 2 {
		t.Errorf("expected 2 deps, got %d", len(b.Deps))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if b.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoadBaseline_Missing(t *testing.T) {
	_, err := gomod.LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadBaseline_Invalid(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := gomod.LoadBaseline(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
