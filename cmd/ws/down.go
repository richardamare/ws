package main

import (
	"fmt"

	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/richardamare/ws/internal/state"
	"github.com/spf13/cobra"
)

func newDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down [project]",
		Short: "Close a project's cmux workspace",
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
			st, err := state.NewStore()
			if err != nil {
				return err
			}
			ref := st.WorkspaceRef(name)
			if ref == "" {
				return fmt.Errorf("no open workspace recorded for %q", name)
			}
			if err := cmuxSvc().Close(bg(), ref); err != nil {
				return err
			}
			_ = st.Clear(name)
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "closed", Value: ref},
			})
		},
	}
}
