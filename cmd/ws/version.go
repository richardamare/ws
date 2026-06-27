package main

import (
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

// version is overridable at build time: -ldflags "-X main.version=..."
var version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the ws version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "version", Value: version},
			})
		},
	}
}
