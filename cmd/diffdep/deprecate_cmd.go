package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runDeprecate(args []string) {
	fs := flag.NewFlagSet("deprecate", flag.ExitOnError)
	branch := fs.String("branch", "", "branch to check (default: current)")
	saveDir := fs.String("save", "", "directory to save the deprecation report (optional)")
	_ = fs.Parse(args)

	client, err := git.NewClient(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *branch == "" {
		cur, err := client.CurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error resolving current branch: %v\n", err)
			os.Exit(1)
		}
		*branch = cur
	}

	loader := &git.DepsLoader{Client: client}
	deps, err := loader.LoadDepsAtBranch(*branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading deps: %v\n", err)
		os.Exit(1)
	}

	report := gomod.CheckDeprecations(deps, *branch)
	fmt.Print(gomod.FormatDeprecationReport(report))

	if *saveDir != "" {
		if err := gomod.SaveDeprecationReport(*saveDir, report); err != nil {
			fmt.Fprintf(os.Stderr, "error saving report: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "report saved to %s\n", *saveDir)
	}

	if len(report.Entries) > 0 {
		os.Exit(1)
	}
}
