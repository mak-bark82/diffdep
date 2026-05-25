package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	branch := fs.String("branch", "", "branch to verify dependencies for")
	outDir := fs.String("out", ".diffdep", "directory to save the verify report")
	save := fs.Bool("save", false, "persist the verify report to disk")
	_ = fs.Parse(args)

	if *branch == "" {
		fmt.Fprintln(os.Stderr, "verify: --branch is required")
		os.Exit(1)
	}

	client, err := git.NewClient(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify: git client: %v\n", err)
		os.Exit(1)
	}

	deps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify: load deps: %v\n", err)
		os.Exit(1)
	}

	// Build a checksum map from go.sum entries (h1: prefix expected).
	checksums := make(map[string]string, len(deps))
	for _, d := range deps {
		if d.Hash != "" {
			checksums[d.Module+"@"+d.Version] = d.Hash
		}
	}

	report := gomod.VerifyDeps(*branch, deps, checksums)
	fmt.Print(gomod.FormatVerifyReport(report))

	if *save {
		if err := gomod.SaveVerifyReport(*outDir, report); err != nil {
			fmt.Fprintf(os.Stderr, "verify: save: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("verify: report saved to %s\n", *outDir)
	}

	if report.Failed > 0 {
		os.Exit(2)
	}
}
