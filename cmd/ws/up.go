package main

import (
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/richardamare/ws/internal/state"
	"github.com/richardamare/ws/internal/workspace"
	"github.com/spf13/cobra"
)

func newUpCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
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

			if dryRun {
				plan := workspace.PlanFor(p)
				return output.Table(cmd.OutOrStdout(), resolveFormat(),
					[]string{"step", "detail"}, plan.Rows())
			}

			azureState := "n/a"
			if p.Azure != nil {
				if err := azureSvc().Login(bg(), p.Azure); err != nil {
					return err
				}
				azureState = "logged-in (Reader)"
			}

			// Background setup (e.g. `docker compose up -d`) before tabs open.
			if err := runScript(bg(), p, p.Setup, flagJSON); err != nil {
				return err
			}

			// Persist the durable cmux.json template so a crash/close can restore
			// the tabs without ws running (ADR-0003). Best-effort: don't fail `up`.
			_, _ = applyTemplate(p, false)

			ref, err := cmuxSvc().Open(bg(), p)
			if err != nil {
				return err
			}
			if st, err := state.NewStore(); err == nil {
				_ = st.SetWorkspaceRef(name, ref)
			}

			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "azure", Value: azureState},
				{Key: "workspace", Value: ref},
				{Key: "tabs", Value: itoa(len(p.Tabs))},
				{Key: "setup", Value: itoa(len(p.Setup))},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show the plan without executing")
	return cmd
}
