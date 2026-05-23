package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewIgnoreList_ShouldIgnore(t *testing.T) {
	il := NewIgnoreList([]string{"github.com/foo", "golang.org/x"})

	cases := []struct {
		module string
		want   bool
	}{
		{"github.com/foo/bar", true},
		{"github.com/foo", true},
		{"golang.org/x/net", true},
		{"github.com/bar/baz", false},
		{"example.com/pkg", false},
	}
	for _, tc := range cases {
		got := il.ShouldIgnore(tc.module)
		if got != tc.want {
			t.Errorf("ShouldIgnore(%q) = %v, want %v", tc.module, got, tc.want)
		}
	}
}

func TestLoadIgnoreFile_Missing(t *testing.T) {
	il, err := LoadIgnoreFile("/nonexistent/.diffdepignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if il.ShouldIgnore("github.com/anything") {
		t.Error("empty ignore list should not ignore anything")
	}
}

func TestLoadIgnoreFile_Parses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".diffdepignore")
	content := "# comment\ngithub.com/ignore-me\n\ngolang.org/x\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	il, err := LoadIgnoreFile(path)
	if err != nil {
		t.Fatalf("LoadIgnoreFile error: %v", err)
	}
	if !il.ShouldIgnore("github.com/ignore-me/pkg") {
		t.Error("expected github.com/ignore-me/pkg to be ignored")
	}
	if !il.ShouldIgnore("golang.org/x/text") {
		t.Error("expected golang.org/x/text to be ignored")
	}
	if il.ShouldIgnore("github.com/keep-me") {
		t.Error("expected github.com/keep-me to not be ignored")
	}
}

func TestApplyIgnore(t *testing.T) {
	diff := DiffResult{
		{Module: "github.com/foo/bar", OldVersion: "v1.0.0", NewVersion: "v2.0.0", ChangeType: ChangeUpdated},
		{Module: "github.com/keep/this", OldVersion: "v1.0.0", NewVersion: "v1.1.0", ChangeType: ChangeUpdated},
		{Module: "golang.org/x/net", OldVersion: "", NewVersion: "v0.1.0", ChangeType: ChangeAdded},
	}
	il := NewIgnoreList([]string{"github.com/foo", "golang.org/x"})
	result := ApplyIgnore(diff, il)

	if len(result) != 1 {
		t.Fatalf("expected 1 entry after ignore, got %d", len(result))
	}
	if result[0].Module != "github.com/keep/this" {
		t.Errorf("unexpected module %q", result[0].Module)
	}
}
