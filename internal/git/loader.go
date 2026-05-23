package git

import (
	"fmt"

	"github.com/yourorg/diffdep/internal/gomod"
)

const (
	GoSumFile = "go.sum"
	GoModFile = "go.mod"
)

// LoadDepsAtBranch reads and parses go.sum (preferred) or go.mod from the
// given branch, returning a map of module path -> version.
func (c *Client) LoadDepsAtBranch(branch string) (map[string]string, error) {
	data, err := c.ReadFileAtBranch(branch, GoSumFile)
	if err == nil {
		deps, parseErr := gomod.ParseGoSum(data)
		if parseErr != nil {
			return nil, fmt.Errorf("parse go.sum at %s: %w", branch, parseErr)
		}
		return gomod.DepsToMap(deps), nil
	}

	// Fallback to go.mod if go.sum is not present.
	data, err = c.ReadFileAtBranch(branch, GoModFile)
	if err != nil {
		return nil, fmt.Errorf("read go.mod at %s: %w", branch, err)
	}
	deps, parseErr := gomod.ParseGoMod(data)
	if parseErr != nil {
		return nil, fmt.Errorf("parse go.mod at %s: %w", branch, parseErr)
	}
	return gomod.DepsToMap(deps), nil
}
