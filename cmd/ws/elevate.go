package main

import (
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

// elevateCommand is the shell run in the elevated tab: drop the scoped config
// dir so az uses the personal admin login, then sign in. Write/Terraform happen
// here, deliberately and human-driven. See docs/security/README.md.
const elevateCommand = `unset AZURE_CONFIG_DIR; echo "ELEVATED — personal admin. Run terraform/az writes deliberately."; az login`

func newElevateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "elevate",
		Short: "Open a marked personal-admin tab for write/Terraform (never uses the Reader SP)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmuxSvc().NewTerminal(bg(), elevateCommand); err != nil {
				return err
			}
			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "elevated", Value: "opened personal-admin tab"},
				{Key: "note", Value: "uses your personal az login, not the Reader SP"},
			})
		},
	}
}
