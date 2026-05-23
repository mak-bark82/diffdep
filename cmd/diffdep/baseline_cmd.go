package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/diffdep/internal/git"
	"github.com/user/diffdep/internal/gomod"
)

// runBaseline handles the "baseline" sub-command:
//   diffdep baseline -branch <branch> -out <file>
//   diffdep baseline -compare <file>
func runBaseline(args []string) error {
	fs := flag.NewFlagSet("baseline", flag.ContinueOnError)
	branch := fs.String("branch", "", "branch to snapshot (default: current)")
	out := fs.String("out", "diffdep-baseline.json", "output file for snapshot")
	compare := fs.String("compare", "", "compare current HEAD against this baseline file")
	repo := fs.String("repo", ".", "path to git repository")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(*repo)
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	if *compare != "" {
		current, err := git.LoadDepsAtBranch(client, "HEAD")
		if err != nil {
			return fmt.Errorf("load current deps: %w", err)
		}
		cmp, err := gomod.CompareAgainstBaseline(*compare, current)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Baseline branch : %s\n", cmp.Baseline.Branch)
		fmt.Fprintf(os.Stdout, "Baseline created: %s\n", cmp.Baseline.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
		if err := gomod.WriteReport(os.Stdout, cmp.Diff, "text"); err != nil {
			return err
		}
		if cmp.HasBreakingChanges() {
			fmt.Fprintln(os.Stderr, "WARNING: breaking version changes detected")
			os.Exit(1)
		}
		return nil
	}

	// Snapshot mode
	if *branch == "" {
		b, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
		*branch = b
	}
	deps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		return fmt.Errorf("load deps: %w", err)
	}
	if err := gomod.SaveBaseline(*out, *branch, deps); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Baseline saved to %s (%d deps from %s)\n", *out, len(deps), *branch)
	return nil
}
