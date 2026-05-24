package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runPrune(args []string) {
	fs := flag.NewFlagSet("prune", flag.ExitOnError)
	base := fs.String("base", "main", "base branch to compare from")
	head := fs.String("head", "", "head branch to compare to (default: current branch)")
	useGoSum := fs.Bool("gosum", false, "parse go.sum instead of go.mod")
	fs.Parse(args)

	client, err := git.NewClient(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising git client: %v\n", err)
		os.Exit(1)
	}

	if *head == "" {
		current, err := client.CurrentBranch()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error detecting current branch: %v\n", err)
			os.Exit(1)
		}
		*head = current
	}

	loader := git.NewLoader(client)

	baseDeps, err := loader.LoadDepsAtBranch(*base, *useGoSum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base deps: %v\n", err)
		os.Exit(1)
	}

	headDeps, err := loader.LoadDepsAtBranch(*head, *useGoSum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading head deps: %v\n", err)
		os.Exit(1)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)
	result := gomod.AnalyzePrune(baseDeps, headDeps, diff)

	fmt.Print(gomod.FormatPruneResult(result))

	if len(result.Suggestions) > 0 {
		os.Exit(1)
	}
}
