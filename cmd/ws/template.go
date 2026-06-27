package main

import (
	"github.com/richardamare/ws/internal/cmux"
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

func newTemplateCmd() *cobra.Command {
	var noReload bool
	cmd := &cobra.Command{
		Use:   "template [project]",
		Short: "Write the project's cmux.json workspace template (durable restore)",
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
			path, err := applyTemplate(p, !noReload)
			if err != nil {
				return err
			}
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "cmux_config", Value: path},
				{Key: "reloaded", Value: boolStr(!noReload)},
			})
		},
	}
	cmd.Flags().BoolVar(&noReload, "no-reload", false, "write the template but don't run cmux reload-config")
	return cmd
}

// applyTemplate upserts the project's workspace template into cmux.json and,
// when reload is true, validates and reloads cmux. Returns the config path.
func applyTemplate(p *config.Project, reload bool) (string, error) {
	path, err := cmux.DefaultConfigPath()
	if err != nil {
		return "", err
	}
	if err := cmux.UpsertCommand(path, cmux.BuildCommand(p)); err != nil {
		return "", err
	}
	if reload {
		svc := cmuxSvc()
		if err := svc.ValidateConfig(bg()); err != nil {
			return path, err
		}
		if err := svc.ReloadConfig(bg()); err != nil {
			return path, err
		}
	}
	return path, nil
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
