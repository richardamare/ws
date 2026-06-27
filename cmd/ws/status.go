package main

import (
	"github.com/richardamare/ws/internal/account"
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/richardamare/ws/internal/state"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [project]",
		Short: "Show a project's config and Azure login state",
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

			rg, azureState := "n/a", "n/a"
			if p.Azure != nil {
				rg = p.Azure.ResourceGroup
				if id, err := azureSvc().Status(bg(), p.Azure); err == nil && id.User == p.Azure.SPAppID {
					azureState = "logged-in (Reader)"
				} else {
					azureState = "logged-out"
				}
			}

			ref := ""
			if st, err := state.NewStore(); err == nil {
				ref = st.WorkspaceRef(name)
			}

			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "name", Value: p.Name},
				{Key: "cwd", Value: p.Cwd},
				{Key: "account", Value: account.FromPath(p.Cwd).Name},
				{Key: "rg", Value: rg},
				{Key: "azure", Value: azureState},
				{Key: "workspace", Value: ref},
				{Key: "tabs", Value: itoa(len(p.Tabs))},
				{Key: "sessions", Value: itoa(len(p.Sessions))},
			})
		},
	}
}
