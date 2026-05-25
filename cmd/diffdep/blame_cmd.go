package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/diffdep/internal/git"
	"github.com/example/diffdep/internal/gomod"
)

func runBlame(args []string) error {
	fs := flag.NewFlagSet("blame", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch to compare against")
	head := fs.String("head", "", "head branch (defaults to current branch)")
	repo := fs.String("repo", ".", "path to git repository")
	file := fs.String("file", "go.sum", "dependency file: go.sum or go.mod")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(*repo)
	if err != nil {
		return fmt.Errorf("blame: failed to open repo: %w", err)
	}

	if *head == "" {
		current, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("blame: could not determine current branch: %w", err)
		}
		*head = current
	}

	loader := git.NewDepsLoader(client, *file)

	baseDeps, err := loader.LoadDepsAtBranch(*base)
	if err != nil {
		return fmt.Errorf("blame: loading base branch %q: %w", *base, err)
	}

	headDeps, err := loader.LoadDepsAtBranch(*head)
	if err != nil {
		return fmt.Errorf("blame: loading head branch %q: %w", *head, err)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)
	report := gomod.NewBlameReport(*head, diff)
	fmt.Fprint(os.Stdout, gomod.FormatBlameReport(report))
	return nil
}
