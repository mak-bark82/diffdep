package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/diffdep/internal/git"
	"github.com/yourorg/diffdep/internal/gomod"
)

func runPin(args []string) error {
	fs := flag.NewFlagSet("pin", flag.ContinueOnError)
	base := fs.String("base", "main", "Base branch to compare against")
	head := fs.String("head", "", "Head branch (defaults to current)")
	pinFile := fs.String("pins", ".diffdep-pins.json", "Path to pin list JSON file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := git.NewClient(".")
	if err != nil {
		return fmt.Errorf("git client: %w", err)
	}

	if *head == "" {
		current, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("current branch: %w", err)
		}
		*head = current
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

	data, err := os.ReadFile(*pinFile)
	if err != nil {
		return fmt.Errorf("read pin file %s: %w", *pinFile, err)
	}

	var pl gomod.PinList
	if err := json.Unmarshal(data, &pl); err != nil {
		return fmt.Errorf("parse pin file: %w", err)
	}

	violations := pl.CheckViolations(diff)
	fmt.Print(gomod.FormatPinViolations(violations))

	if len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
