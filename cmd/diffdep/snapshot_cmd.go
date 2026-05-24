package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/diffdep/internal/git"
	"github.com/your-org/diffdep/internal/gomod"
)

// runSnapshot handles the "snapshot" subcommand.
// Usage: diffdep snapshot -branch <branch> -dir <dir> [-repo <repo>]
func runSnapshot(args []string) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	branch := fs.String("branch", "", "branch to snapshot (required)")
	dir := fs.String("dir", ".diffdep/snapshots", "directory to store snapshots")
	repo := fs.String("repo", ".", "path to git repository")
	diff := fs.String("diff", "", "compare against this branch and print diff")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *branch == "" {
		return fmt.Errorf("snapshot: -branch is required")
	}

	client, err := git.NewClient(*repo)
	if err != nil {
		return fmt.Errorf("snapshot: git client: %w", err)
	}

	deps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		return fmt.Errorf("snapshot: load deps: %w", err)
	}

	s := gomod.NewSnapshot(*branch, deps)
	if err := gomod.SaveSnapshot(*dir, s); err != nil {
		return fmt.Errorf("snapshot: save: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Snapshot saved: %s\n", gomod.SnapshotSummary(s))

	if *diff != "" {
		base, err := gomod.LoadSnapshot(*dir, *diff)
		if err != nil {
			return fmt.Errorf("snapshot: load base %q: %w", *diff, err)
		}
		result := gomod.DiffSnapshot(base, s)
		summary := gomod.Summarize(result)
		fmt.Fprintln(os.Stdout, summary.String())
		if err := gomod.WriteReport(os.Stdout, result, "text"); err != nil {
			return fmt.Errorf("snapshot: write report: %w", err)
		}
	}
	return nil
}
