package main

import (
	"fmt"

	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

func newRotateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rotate [project]",
		Short: "Issue a fresh certificate for the project's Reader SP",
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

			sp, err := azureSvc().RotateCert(bg(), p.Azure.SPAppID)
			if err != nil {
				return err
			}
			if err := copyFile(sp.CertFile, config.ExpandHome(p.Azure.Cert)); err != nil {
				return fmt.Errorf("store rotated cert: %w", err)
			}
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "rotated_cert", Value: p.Azure.Cert},
			})
		},
	}
}
