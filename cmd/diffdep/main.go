package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func main() {
	base := flag.String("base", "main", "base branch to compare from")
	head := flag.String("head", "", "head branch to compare to (default: current branch)")
	format := flag.String("format", "text", "output format: text or markdown")
	majorOnly := flag.Bool("major-only", false, "only show major version changes")
	prefix := flag.String("prefix", "", "filter dependencies by module prefix")
	flag.Parse()

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

	baseDeps, err := git.LoadDepsAtBranch(client, *base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading deps for %s: %v\n", *base, err)
		os.Exit(1)
	}

	headDeps, err := git.LoadDepsAtBranch(client, *head)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading deps for %s: %v\n", *head, err)
		os.Exit(1)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)

	opts := gomod.FilterOptions{
		MajorOnly:     *majorOnly,
		Prefix:        *prefix,
		IncludeAdded:  true,
		IncludeRemoved: true,
	}
	diff = gomod.FilterDiff(diff, opts)

	summary := gomod.Summarize(diff)
	fmt.Fprintf(os.Stderr, "summary: %s\n", summary)

	if err := gomod.WriteReport(os.Stdout, diff, gomod.ReportFormat(*format)); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}
}
