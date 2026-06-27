package main

import (
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List configured projects",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := config.NewStore()
			if err != nil {
				return err
			}
			names, err := store.List()
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(names))
			for _, name := range names {
				p, err := store.Load(name)
				if err != nil {
					rows = append(rows, []string{name, "", "(unreadable)"})
					continue
				}
				rg := ""
				if p.Azure != nil {
					rg = p.Azure.ResourceGroup
				}
				rows = append(rows, []string{p.Name, p.Cwd, rg})
			}
			return output.Table(cmd.OutOrStdout(), resolveFormat(),
				[]string{"name", "cwd", "rg"}, rows)
		},
	}
}
