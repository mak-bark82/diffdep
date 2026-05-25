package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runDrift(args []string) error {
	fs := flag.NewFlagSet("drift", flag.ContinueOnError)
	branch := fs.String("branch", "main", "branch to analyse")
	baselineFile := fs.String("baseline", ".diffdep/baseline.json", "baseline file path")
	outDir := fs.String("out", ".diffdep", "directory to save drift report")
	sinceDays := fs.Int("since", 30, "days since baseline was captured")
	save := fs.Bool("save", false, "persist the drift report to disk")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("drift: git client: %w", err)
	}

	currentDeps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		return fmt.Errorf("drift: load deps: %w", err)
	}

	baseline, err := gomod.LoadBaseline(*baselineFile)
	if err != nil {
		return fmt.Errorf("drift: load baseline: %w", err)
	}

	since := time.Now().Add(-time.Duration(*sinceDays) * 24 * time.Hour)
	report := gomod.AnalyzeDrift(*branch, baseline, currentDeps, since)

	fmt.Print(gomod.FormatDriftReport(report))

	if *save {
		if err := gomod.SaveDriftReport(*outDir, report); err != nil {
			return fmt.Errorf("drift: save: %w", err)
		}
		fmt.Fprintf(os.Stderr, "drift report saved to %s/drift.json\n", *outDir)
	}

	return nil
}
