package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runOutdated(args []string) error {
	fs := flag.NewFlagSet("outdated", flag.ContinueOnError)
	branch := fs.String("branch", "", "Branch to check (default: current branch)")
	repo := fs.String("repo", ".", "Path to git repository")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(*repo)
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	if *branch == "" {
		cur, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
		*branch = cur
	}

	deps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		return fmt.Errorf("load deps: %w", err)
	}

	// Default resolver: treats the version in go.sum as canonical (stub for CLI).
	// A real implementation would query proxy.golang.org.
	resolver := func(module string) (string, error) {
		for _, d := range deps {
			if d.Module == module {
				return d.Version, nil
			}
		}
		return "", fmt.Errorf("module %s not found", module)
	}

	report, err := gomod.CheckOutdated(*branch, deps, resolver)
	if err != nil {
		return fmt.Errorf("check outdated: %w", err)
	}

	fmt.Fprint(os.Stdout, gomod.FormatOutdatedReport(report))
	return nil
}
