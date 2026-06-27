package main

import (
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/richardamare/ws/internal/workspace"
	"github.com/spf13/cobra"
)

func newUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up [project]",
		Short: "Start working on a project: scoped Azure login + cmux workspace",
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
			// The plan is computed now; execution (cmux + az) lands in follow-up
			// commits. `up` already reports exactly what it will do.
			plan := workspace.PlanFor(p)
			return output.Table(cmd.OutOrStdout(), resolveFormat(),
				[]string{"step", "detail"}, plan.Rows())
		},
	}
}
