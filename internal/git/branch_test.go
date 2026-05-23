package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/yourorg/diffdep/internal/git"
)

func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}

	run("init")
	run("checkout", "-b", "main")

	gosum := filepath.Join(dir, "go.sum")
	if err := os.WriteFile(gosum, []byte("github.com/pkg/errors v0.9.1 h1:abc\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")
	return dir
}

func TestCurrentBranch(t *testing.T) {
	dir := initTestRepo(t)
	c := git.NewClient(dir)
	branch, err := c.CurrentBranch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "main" {
		t.Errorf("expected main, got %q", branch)
	}
}

func TestReadFileAtBranch(t *testing.T) {
	dir := initTestRepo(t)
	c := git.NewClient(dir)
	data, err := c.ReadFileAtBranch("main", "go.sum")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty file content")
	}
}

func TestListBranches(t *testing.T) {
	dir := initTestRepo(t)
	c := git.NewClient(dir)
	branches, err := c.ListBranches()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) == 0 {
		t.Error("expected at least one branch")
	}
}
