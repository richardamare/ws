package main

import (
	"fmt"

	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/richardamare/ws/internal/session"
	"github.com/spf13/cobra"
)

func newSaveCmd() *cobra.Command {
	var note string
	cmd := &cobra.Command{
		Use:   "save <project> <label>",
		Short: "Bookmark the current Claude session under a label",
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
			if len(args) < 2 {
				return fmt.Errorf("label required: ws save %s <label>", name)
			}
			label := args[1]

			id, err := cmuxSvc().ResumeID(bg())
			if err != nil {
				return fmt.Errorf("read current session id from cmux: %w", err)
			}
			if id == "" {
				return fmt.Errorf("no resumable agent session in the focused cmux surface")
			}

			p, err := store.Load(name)
			if err != nil {
				return err
			}
			p.Sessions = session.Upsert(p.Sessions, config.Bookmark{Label: label, ID: id, Note: note})
			if err := store.Save(p); err != nil {
				return err
			}
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "label", Value: label},
				{Key: "id", Value: id},
			})
		},
	}
	cmd.Flags().StringVar(&note, "note", "", "optional note for the bookmark")
	return cmd
}
