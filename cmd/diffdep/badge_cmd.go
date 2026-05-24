package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runBadge(args []string) error {
	fs := flag.NewFlagSet("badge", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch")
	head := fs.String("head", "", "head branch (default: current)")
	style := fs.String("style", "flat", "badge style: flat, flat-square, plastic")
	format := fs.String("format", "markdown", "output format: markdown, svg, json")
	outDir := fs.String("out", "", "directory to save badge.json (optional)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	if *head == "" {
		cur, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
		*head = cur
	}

	baseDeps, err := git.LoadDepsAtBranch(client, *base)
	if err != nil {
		return fmt.Errorf("load base deps: %w", err)
	}
	headDeps, err := git.LoadDepsAtBranch(client, *head)
	if err != nil {
		return fmt.Errorf("load head deps: %w", err)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)
	score := gomod.ScoreDiff(diff, gomod.DefaultWeights)
	badge := gomod.GenerateBadge(score, gomod.BadgeStyle(*style))

	if *outDir != "" {
		if err := gomod.SaveBadge(*outDir, badge); err != nil {
			return fmt.Errorf("save badge: %w", err)
		}
	}

	switch *format {
	case "markdown":
		fmt.Println(gomod.FormatBadgeMarkdown(badge))
	case "svg":
		fmt.Println(gomod.FormatBadgeSVG(badge))
	case "json":
		fmt.Printf(`{"label":%q,"message":%q,"color":%q,"style":%q}`+"\n",
			badge.Label, badge.Message, badge.Color, badge.Style)
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q, using markdown\n", *format)
		fmt.Println(gomod.FormatBadgeMarkdown(badge))
	}

	return nil
}
