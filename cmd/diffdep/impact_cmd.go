package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runImpact(args []string) error {
	fs := flag.NewFlagSet("impact", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch to compare against")
	target := fs.String("target", "", "target branch (defaults to current branch)")
	format := fs.String("format", "text", "output format: text")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	if *target == "" {
		cur, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
		*target = cur
	}

	loader := &git.DepLoader{Client: client}

	baseDeps, err := loader.LoadDepsAtBranch(*base)
	if err != nil {
		return fmt.Errorf("load base deps: %w", err)
	}

	targetDeps, err := loader.LoadDepsAtBranch(*target)
	if err != nil {
		return fmt.Errorf("load target deps: %w", err)
	}

	diff := gomod.DiffDependencies(baseDeps, targetDeps)
	report := gomod.AssessImpact(*target, diff)

	switch *format {
	case "text":
		fmt.Print(gomod.FormatImpactReport(report))
	default:
		return fmt.Errorf("unsupported format: %s", *format)
	}

	hasHigh := false
	for _, e := range report.Entries {
		if e.Level == gomod.ImpactHigh || e.Level == gomod.ImpactCritical {
			hasHigh = true
			break
		}
	}
	if hasHigh {
		os.Exit(1)
	}
	return nil
}
