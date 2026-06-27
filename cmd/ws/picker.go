package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/richardamare/ws/internal/config"
)

// resolveProjectName returns the project name from args, or — on an interactive
// TTY — prompts the user to pick one. Under --json/non-TTY a missing name is an
// error (no prompts). See docs/patterns/cli.md.
func resolveProjectName(store *config.Store, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	if !interactive() {
		return "", fmt.Errorf("project name required")
	}
	names, err := store.List()
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no projects configured; run `ws new` first")
	}
	var choice string
	opts := make([]huh.Option[string], 0, len(names))
	for _, n := range names {
		opts = append(opts, huh.NewOption(n, n))
	}
	if err := huh.NewSelect[string]().
		Title("Which project?").
		Options(opts...).
		Value(&choice).
		Run(); err != nil {
		return "", err
	}
	return choice, nil
}
