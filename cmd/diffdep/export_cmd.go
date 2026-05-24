package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

// runExport implements the "export" subcommand. It computes the dependency diff
// between a base branch and a head branch, then writes the result to stdout or
// a file in the requested format (json or csv).
func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch")
	head := fs.String("head", "", "head branch (default: current)")
	format := fs.String("format", "json", "output format: json or csv")
	output := fs.String("output", "", "output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *format != "json" && *format != "csv" {
		return fmt.Errorf("unsupported format %q: must be \"json\" or \"csv\"", *format)
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	headBranch := *head
	if headBranch == "" {
		headBranch, err = client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
	}

	baseDeps, err := git.LoadDepsAtBranch(client, *base)
	if err != nil {
		return fmt.Errorf("load base deps: %w", err)
	}

	headDeps, err := git.LoadDepsAtBranch(client, headBranch)
	if err != nil {
		return fmt.Errorf("load head deps: %w", err)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if err := gomod.ExportDiff(w, diff, gomod.ExportFormat(*format)); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if *output != "" {
		fmt.Fprintf(os.Stderr, "exported %d entries to %s\n", len(diff), *output)
	}
	return nil
}
