package main

import (
	"os"

	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	flagJSON  bool
	flagPlain bool
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "ws",
		Short:         "Set up a per-project developer workspace (cmux tabs, scoped Azure login, session bookmarks)",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().BoolVar(&flagJSON, "json", false, "strict JSON output; disables interactive prompts")
	root.PersistentFlags().BoolVar(&flagPlain, "plain", false, "structured text output (auto on non-TTY)")

	root.AddCommand(
		newVersionCmd(),
		newNewCmd(),
		newListCmd(),
		newStatusCmd(),
		newUpCmd(),
		newDownCmd(),
		newAuthCmd(),
		newRotateCmd(),
		newElevateCmd(),
		newSessionsCmd(),
		newSaveCmd(),
		newResumeCmd(),
		newRmCmd(),
	)
	return root
}

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

// resolveFormat picks the output format from the global flags and whether
// stdout is a terminal.
func resolveFormat() output.Format {
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	return output.Detect(flagJSON, flagPlain, isTTY)
}

// interactive reports whether interactive pickers (huh) are allowed: only on a
// real TTY and never under --json. Missing input otherwise is an error.
func interactive() bool {
	return !flagJSON && term.IsTerminal(int(os.Stdin.Fd()))
}
