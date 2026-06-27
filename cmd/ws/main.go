// Command ws sets up a per-project developer workspace: cmux tabs, a scoped
// Azure login, and Claude Code session bookmarks.
//
// This package is the application interface (the CLI). All business logic lives
// in internal/*; the files here only parse input and render results.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "ws:", err)
		os.Exit(1)
	}
}
