package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/richardamare/ws/internal/state"
	"github.com/spf13/cobra"
)

func newRmCmd() *cobra.Command {
	var purge, yes bool
	cmd := &cobra.Command{
		Use:   "rm <project>",
		Short: "Remove a project's config (and with --purge its cert)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			store, err := config.NewStore()
			if err != nil {
				return err
			}
			p, err := store.Load(name)
			if err != nil {
				return err
			}

			if !yes {
				if !interactive() {
					return fmt.Errorf("refusing to remove %q without --yes", name)
				}
				var ok bool
				if err := huh.NewConfirm().Title(fmt.Sprintf("Remove project %q?", name)).Value(&ok).Run(); err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("aborted")
				}
			}

			if err := store.Remove(name); err != nil {
				return err
			}
			if st, err := state.NewStore(); err == nil {
				_ = st.Clear(name)
			}

			spHint := ""
			if purge && p.Azure != nil {
				_ = os.Remove(config.ExpandHome(p.Azure.Cert))
				// SP deletion is deliberate and destructive — never automated.
				spHint = fmt.Sprintf("az ad sp delete --id %s", p.Azure.SPAppID)
			}

			fields := []output.KV{{Key: "removed", Value: name}}
			if spHint != "" {
				fields = append(fields, output.KV{Key: "delete_sp_manually", Value: spHint})
			}
			return output.Record(cmd.OutOrStdout(), resolveFormat(), fields)
		},
	}
	cmd.Flags().BoolVar(&purge, "purge", false, "also delete the local cert and print the SP-deletion command")
	cmd.Flags().BoolVar(&yes, "yes", false, "skip confirmation")
	return cmd
}
