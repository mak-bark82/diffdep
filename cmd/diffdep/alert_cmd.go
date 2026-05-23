package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/diffdep/internal/git"
	"github.com/user/diffdep/internal/gomod"
)

// runAlert compares two branches and prints dependency alerts.
func runAlert(args []string) error {
	fs := flag.NewFlagSet("alert", flag.ContinueOnError)
	base := fs.String("base", "main", "base branch")
	head := fs.String("head", "", "head branch (default: current)")
	majorOnly := fs.Bool("major-only", false, "alert on major version changes only")
	saveFile := fs.String("save", "", "save alerts to JSON file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	headBranch := *head
	if headBranch == "" {
		headBranch, err = client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
	}

	baseDeps, err := git.LoadDepsAtBranch(client, *base)
	if err != nil {
		return fmt.Errorf("load base deps: %w", err)
	}
	headDeps, err := git.LoadDepsAtBranch(client, headBranch)
	if err != nil {
		return fmt.Errorf("load head deps: %w", err)
	}

	diff := gomod.DiffDependencies(baseDeps, headDeps)

	cfg := gomod.DefaultAlertConfig()
	if *majorOnly {
		cfg = gomod.AlertConfig{OnMajor: true}
	}

	alerts := gomod.GenerateAlerts(diff, cfg)
	fmt.Println(gomod.FormatAlerts(alerts))

	if *saveFile != "" {
		if err := gomod.SaveAlerts(*saveFile, alerts); err != nil {
			return fmt.Errorf("save alerts: %w", err)
		}
		fmt.Fprintf(os.Stderr, "alerts saved to %s\n", *saveFile)
	}

	if len(alerts) > 0 {
		for _, a := range alerts {
			if a.Level == gomod.AlertCritical {
				os.Exit(2)
			}
		}
		os.Exit(1)
	}
	return nil
}
