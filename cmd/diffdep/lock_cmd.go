package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/diffdep/internal/git"
	"github.com/user/diffdep/internal/gomod"
)

func runLock(args []string) {
	fs := flag.NewFlagSet("lock", flag.ExitOnError)
	branch := fs.String("branch", "main", "branch to lock dependencies from")
	dir := fs.String("dir", ".diffdep", "directory to store lock file")
	check := fs.Bool("check", false, "check current branch against lock file instead of writing")
	_ = fs.Parse(args)

	client, err := git.NewClient(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "lock: git client: %v\n", err)
		os.Exit(1)
	}

	if *check {
		lf, err := gomod.LoadLockFile(*dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lock: load: %v\n", err)
			os.Exit(1)
		}
		if lf == nil {
			fmt.Fprintln(os.Stderr, "lock: no lock file found; run without --check to create one")
			os.Exit(1)
		}
		current, err := git.LoadDepsAtBranch(client, "HEAD")
		if err != nil {
			fmt.Fprintf(os.Stderr, "lock: load deps: %v\n", err)
			os.Exit(1)
		}
		locked, err := gomod.LoadLockFile(*dir)
		if err != nil || locked == nil {
			fmt.Fprintf(os.Stderr, "lock: reload: %v\n", err)
			os.Exit(1)
		}
		baseDeps := make([]gomod.Dependency, 0, len(locked.Entries))
		for _, e := range locked.Entries {
			baseDeps = append(baseDeps, gomod.Dependency{Module: e.Module, Version: e.Version})
		}
		diff := gomod.DiffDependencies(baseDeps, current)
		violations := gomod.CheckLock(lf, diff)
		fmt.Print(gomod.FormatLockViolations(violations))
		if len(violations) > 0 {
			os.Exit(2)
		}
		return
	}

	deps, err := git.LoadDepsAtBranch(client, *branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lock: load deps at %s: %v\n", *branch, err)
		os.Exit(1)
	}
	lf := gomod.NewLockFile(deps)
	if err := gomod.SaveLockFile(*dir, lf); err != nil {
		fmt.Fprintf(os.Stderr, "lock: save: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Lock file written to %s with %d entries from branch %s\n", *dir, len(lf.Entries), *branch)
}
