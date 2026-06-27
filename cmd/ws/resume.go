package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/session"
	"github.com/spf13/cobra"
)

func newResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume <project> [label]",
		Short: "Resume a bookmarked Claude session",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := config.NewStore()
			if err != nil {
				return err
			}
			name, err := resolveProjectName(store, args)
			if err != nil {
				return err
			}
			p, err := store.Load(name)
			if err != nil {
				return err
			}
			if len(p.Sessions) == 0 {
				return fmt.Errorf("no bookmarks for %q; use `ws save`", name)
			}

			label := ""
			if len(args) >= 2 {
				label = args[1]
			} else {
				if !interactive() {
					return fmt.Errorf("label required in non-interactive mode")
				}
				if label, err = pickBookmark(p.Sessions); err != nil {
					return err
				}
			}

			b, ok := session.Find(p.Sessions, label)
			if !ok {
				return fmt.Errorf("no bookmark %q in %q", label, name)
			}
			resumeArgs, err := session.ResumeArgs(b)
			if err != nil {
				return err
			}
			return stdio(bg(), "claude", resumeArgs...)
		},
	}
}

func pickBookmark(bookmarks []config.Bookmark) (string, error) {
	opts := make([]huh.Option[string], 0, len(bookmarks))
	for _, b := range bookmarks {
		label := b.Label
		if b.Note != "" {
			label = b.Label + " — " + b.Note
		}
		opts = append(opts, huh.NewOption(label, b.Label))
	}
	var choice string
	if err := huh.NewSelect[string]().Title("Resume which session?").Options(opts...).Value(&choice).Run(); err != nil {
		return "", err
	}
	return choice, nil
}
