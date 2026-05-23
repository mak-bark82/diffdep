package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client provides git operations for a repository.
type Client struct {
	RepoPath string
}

// NewClient creates a new git Client for the given repo path.
func NewClient(repoPath string) *Client {
	return &Client{RepoPath: repoPath}
}

// ReadFileAtBranch returns the contents of a file at a given branch or commit ref.
func (c *Client) ReadFileAtBranch(branch, filePath string) ([]byte, error) {
	ref := fmt.Sprintf("%s:%s", branch, filePath)
	cmd := exec.Command("git", "-C", c.RepoPath, "show", ref)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show %s: %w", ref, err)
	}
	return out, nil
}

// ListBranches returns all local branch names in the repository.
func (c *Client) ListBranches() ([]string, error) {
	cmd := exec.Command("git", "-C", c.RepoPath, "branch", "--format=%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git branch: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var branches []string
	for _, l := range lines {
		if t := strings.TrimSpace(l); t != "" {
			branches = append(branches, t)
		}
	}
	return branches, nil
}

// CurrentBranch returns the name of the currently checked-out branch.
func (c *Client) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "-C", c.RepoPath, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
