package main

import (
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

func newSessionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sessions [project]",
		Short: "List a project's curated Claude session bookmarks",
		Args:  cobra.MaximumNArgs(1),
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
			rows := make([][]string, 0, len(p.Sessions))
			for _, b := range p.Sessions {
				rows = append(rows, []string{b.Label, b.ID, b.Note})
			}
			return output.Table(cmd.OutOrStdout(), resolveFormat(),
				[]string{"label", "id", "note"}, rows)
		},
	}
}
