package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runOverlap(args []string) error {
	fs := flag.NewFlagSet("overlap", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch")
	head := fs.String("head", "", "head branch to compare (required)")
	repo := fs.String("repo", ".", "path to git repository")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *head == "" {
		fs.Usage()
		return fmt.Errorf("--head is required")
	}

	client, err := git.NewClient(*repo)
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	loader := git.NewLoader(client)

	baseDeps, err := loader.LoadDepsAtBranch(*base)
	if err != nil {
		return fmt.Errorf("loading base branch %q: %w", *base, err)
	}

	headDeps, err := loader.LoadDepsAtBranch(*head)
	if err != nil {
		return fmt.Errorf("loading head branch %q: %w", *head, err)
	}

	report := gomod.AnalyzeOverlap(*head, baseDeps, headDeps)
	fmt.Fprint(os.Stdout, gomod.FormatOverlapReport(report))
	return nil
}
