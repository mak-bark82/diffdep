package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runAnnotate(args []string) {
	fs := flag.NewFlagSet("annotate", flag.ExitOnError)
	base := fs.String("base", "main", "base branch to compare from")
	head := fs.String("head", "", "head branch to compare to (default: current branch)")
	useGoMod := fs.Bool("gomod", false, "parse go.mod instead of go.sum")
	fs.Parse(args) //nolint:errcheck

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

	baseDeps, err := loader.LoadDepsAtBranch(*base, *useGoMod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base deps: %v\n", err)
		os.Exit(1)
	}

	headDeps, err := loader.LoadDepsAtBranch(*head, *useGoMod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading head deps: %v\n", err)
		os.Exit(1)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)
	annotations := gomod.AnnotateDiff(diff)

	fmt.Printf("Annotations for %s → %s:\n\n", *base, *head)
	fmt.Print(gomod.FormatAnnotations(annotations))
}
