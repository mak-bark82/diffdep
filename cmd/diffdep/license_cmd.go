package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runLicense(args []string) error {
	fs := flag.NewFlagSet("license", flag.ContinueOnError)
	branch := fs.String("branch", "", "branch to inspect (default: current)")
	saveDir := fs.String("save", "", "directory to persist the license report")
	highOnly := fs.Bool("high-only", false, "only show high-risk licenses")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("license: git client: %w", err)
	}

	target := *branch
	if target == "" {
		target, err = client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("license: current branch: %w", err)
		}
	}

	deps, err := git.LoadDepsAtBranch(client, target)
	if err != nil {
		return fmt.Errorf("license: load deps: %w", err)
	}

	report := gomod.CheckLicenses(deps)

	if *highOnly {
		filtered := make([]gomod.LicenseEntry, 0)
		for _, e := range report.Entries {
			if e.Risk == "high" {
				filtered = append(filtered, e)
			}
		}
		report.Entries = filtered
	}

	fmt.Fprintln(os.Stdout, gomod.FormatLicenseReport(report))

	if *saveDir != "" {
		if err := gomod.SaveLicenseReport(*saveDir, report); err != nil {
			return fmt.Errorf("license: save: %w", err)
		}
		fmt.Fprintf(os.Stdout, "report saved to %s\n", *saveDir)
	}

	if report.HighRiskCount > 0 {
		return fmt.Errorf("license: %d high-risk license(s) detected", report.HighRiskCount)
	}
	return nil
}
