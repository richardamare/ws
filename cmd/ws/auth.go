package main

import (
	"fmt"

	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth [project]",
		Short: "Log the project's Reader SP into its isolated config dir",
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
			if p.Azure == nil {
				return fmt.Errorf("project %q has no azure config", name)
			}
			if err := azureSvc().Login(bg(), p.Azure); err != nil {
				return err
			}
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "azure", Value: "logged-in (Reader)"},
			})
		},
	}
}
