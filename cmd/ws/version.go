package main

import (
	"runtime/debug"

	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

// version is overridable at build time: -ldflags "-X main.version=..." (set by
// GoReleaser). When unset, it falls back to the module version embedded by
// `go install ...@vX.Y.Z`.
var version = "dev"

func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return version
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the ws version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "version", Value: resolveVersion()},
			})
		},
	}
}
