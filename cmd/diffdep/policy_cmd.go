package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runPolicy(args []string) error {
	fs := flag.NewFlagSet("policy", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch to compare against")
	head := fs.String("head", "", "head branch (defaults to current branch)")
	policyFile := fs.String("policy", ".diffdep-policy.json", "path to policy JSON file")
	useGoSum := fs.Bool("gosum", false, "use go.sum instead of go.mod")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("open repo: %w", err)
	}

	headBranch := *head
	if headBranch == "" {
		headBranch, err = client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
	}

	loader := git.NewLoader(client)
	baseDeps, err := loader.LoadDepsAtBranch(*base, *useGoSum)
	if err != nil {
		return fmt.Errorf("load base deps: %w", err)
	}
	headDeps, err := loader.LoadDepsAtBranch(headBranch, *useGoSum)
	if err != nil {
		return fmt.Errorf("load head deps: %w", err)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)

	policy, err := gomod.LoadPolicy(*policyFile)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}

	violations := gomod.EnforcePolicy(policy, diff)
	if len(violations) == 0 {
		fmt.Println("✓ All policy checks passed.")
		return nil
	}

	fmt.Fprintf(os.Stderr, "✗ Policy violations found (%d):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  %s\n", v)
	}
	return fmt.Errorf("policy enforcement failed with %d violation(s)", len(violations))
}
