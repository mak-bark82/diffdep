package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

// runAudit compares dependencies between two branches and produces an audit
// report. By default it compares the current branch against "main", but both
// the target branch and the base branch can be overridden via flags.
func runAudit(args []string) error {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	branch := fs.String("branch", "", "branch to audit (default: current)")
	base := fs.String("base", "main", "base branch to compare against")
	outDir := fs.String("out", ".diffdep", "directory to save audit report")
	print := fs.Bool("print", false, "print audit report to stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("audit: git client: %w", err)
	}

	targetBranch := *branch
	if targetBranch == "" {
		targetBranch, err = client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("audit: current branch: %w", err)
		}
	}

	if targetBranch == *base {
		return fmt.Errorf("audit: target branch %q is the same as base branch %q; nothing to compare", targetBranch, *base)
	}

	loader := git.NewDepsLoader(client)

	baseDeps, err := loader.LoadDepsAtBranch(*base)
	if err != nil {
		return fmt.Errorf("audit: load base deps: %w", err)
	}

	headDeps, err := loader.LoadDepsAtBranch(targetBranch)
	if err != nil {
		return fmt.Errorf("audit: load head deps: %w", err)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)
	report := gomod.NewAuditReport(targetBranch, diff)

	if *print {
		fmt.Print(gomod.FormatAuditReport(report))
	}

	if err := gomod.SaveAuditReport(*outDir, report); err != nil {
		return fmt.Errorf("audit: save: %w", err)
	}

	fmt.Fprintf(os.Stderr, "audit: report saved to %s/audit.json (%d entries)\n", *outDir, len(report.Entries))
	return nil
}
