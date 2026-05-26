package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runRename(args []string) {
	fs := flag.NewFlagSet("rename", flag.ExitOnError)
	base := fs.String("base", "main", "base branch to compare from")
	head := fs.String("head", "", "head branch to compare to (required)")
	repo := fs.String("repo", ".", "path to git repository")
	output := fs.String("output", "", "optional path to save JSON report")
	fs.Parse(args)

	if *head == "" {
		fmt.Fprintln(os.Stderr, "error: --head branch is required")
		os.Exit(1)
	}

	client, err := git.NewClient(*repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: git client: %v\n", err)
		os.Exit(1)
	}

	baseDeps, err := git.LoadDepsAtBranch(client, *base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: load base deps: %v\n", err)
		os.Exit(1)
	}

	headDeps, err := git.LoadDepsAtBranch(client, *head)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: load head deps: %v\n", err)
		os.Exit(1)
	}

	baseMap := gomod.DepsToMap(baseDeps)
	headMap := gomod.DepsToMap(headDeps)

	baseStrMap := make(map[string]string, len(baseMap))
	for k, v := range baseMap {
		baseStrMap[k] = v.Version
	}
	headStrMap := make(map[string]string, len(headMap))
	for k, v := range headMap {
		headStrMap[k] = v.Version
	}

	report := gomod.DetectRenames(baseStrMap, headStrMap, *head)
	fmt.Print(gomod.FormatRenameReport(report))

	if *output != "" {
		if err := gomod.SaveRenameReport(*output, report); err != nil {
			fmt.Fprintf(os.Stderr, "error: save report: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "report saved to %s\n", *output)
	}
}
