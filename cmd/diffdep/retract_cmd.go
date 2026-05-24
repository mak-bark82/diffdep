package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runRetract(args []string) {
	fs := flag.NewFlagSet("retract", flag.ExitOnError)
	branch := fs.String("branch", "", "branch to check for retracted modules (default: current)")
	useGoSum := fs.Bool("gosum", false, "parse go.sum instead of go.mod")
	fs.Parse(args)

	client, err := git.NewClient(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising git client: %v\n", err)
		os.Exit(1)
	}

	if *branch == "" {
		cur, err := client.CurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting current branch: %v\n", err)
			os.Exit(1)
		}
		*branch = cur
	}

	fileName := "go.mod"
	if *useGoSum {
		fileName = "go.sum"
	}

	raw, err := client.ReadFileAtBranch(*branch, fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s at branch %s: %v\n", fileName, *branch, err)
		os.Exit(1)
	}

	var deps []gomod.Dependency
	if *useGoSum {
		deps, err = gomod.ParseGoSum(raw)
	} else {
		deps, err = gomod.ParseGoMod(raw)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing dependencies: %v\n", err)
		os.Exit(1)
	}

	report := gomod.CheckRetractions(*branch, deps)
	fmt.Print(gomod.FormatRetractReport(report))

	if len(report.Retracted) > 0 {
		os.Exit(2)
	}
}
