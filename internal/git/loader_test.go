package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/yourorg/diffdep/internal/git"
)

func initRepoWithGoSum(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=t@t.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=t@t.com",
		)
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	run("init")
	run("checkout", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "go.sum"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")
	return dir
}

func TestLoadDepsAtBranch_GoSum(t *testing.T) {
	content := "github.com/pkg/errors v0.9.1 h1:fakehash\ngithub.com/stretchr/testify v1.8.0 h1:fakehash\n"
	dir := initRepoWithGoSum(t, content)
	c := git.NewClient(dir)

	deps, err := c.LoadDepsAtBranch("main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := deps["github.com/pkg/errors"]; !ok || v != "v0.9.1" {
		t.Errorf("expected pkg/errors v0.9.1, got %q", v)
	}
	if v, ok := deps["github.com/stretchr/testify"]; !ok || v != "v1.8.0" {
		t.Errorf("expected testify v1.8.0, got %q", v)
	}
}

func TestLoadDepsAtBranch_MissingFile(t *testing.T) {
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=t@t.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=t@t.com",
		)
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	run("init")
	run("checkout", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")

	c := git.NewClient(dir)
	_, err := c.LoadDepsAtBranch("main")
	if err == nil {
		t.Error("expected error when go.sum and go.mod are missing")
	}
}
