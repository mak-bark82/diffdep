package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runNotify(args []string) {
	fs := flag.NewFlagSet("notify", flag.ExitOnError)
	branch := fs.String("branch", "", "branch to report on (required)")
	base := fs.String("base", "main", "base branch to compare against")
	channel := fs.String("channel", "stdout", "notification channel: stdout, slack, webhook")
	format := fs.String("format", "text", "output format: text, markdown")
	minRisk := fs.String("min-risk", "medium", "minimum risk level to notify: low, medium, high, critical")
	_ = fs.Parse(args)

	if *branch == "" {
		fmt.Fprintln(os.Stderr, "error: --branch is required")
		os.Exit(1)
	}

	client := git.NewClient(".")

	baseDeps, err := git.LoadDepsAtBranch(client, *base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base branch %q: %v\n", *base, err)
		os.Exit(1)
	}

	headDeps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading branch %q: %v\n", *branch, err)
		os.Exit(1)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)
	summary := gomod.Summarize(diff)
	score := gomod.ScoreDiff(diff, gomod.DefaultWeights())
	alerts := gomod.GenerateAlerts(diff, gomod.DefaultAlertConfig())

	cfg := gomod.NotifyConfig{
		Channel:    gomod.NotifyChannel(*channel),
		MinRisk:    *minRisk,
		Branch:     *branch,
	}

	payload := gomod.NotifyPayload{
		Branch:  cfg.Branch,
		Summary: summary,
		Score:   score,
		Alerts:  alerts,
	}

	var output string
	switch *format {
	case "markdown":
		output = gomod.FormatNotifyMarkdown(payload)
	default:
		output = gomod.FormatNotifyText(payload)
	}

	fmt.Print(output)
}
