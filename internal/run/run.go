// Package run is the seam between ws and the external tools it drives (az, cmux,
// claude, docker). Services take a Runner so tests can inject a fake instead of
// spawning real processes. See docs/patterns/go.md.
package run

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// Runner executes an external command and returns its stdout.
type Runner interface {
	Run(ctx context.Context, env []string, name string, args ...string) (string, error)
}

// Exec is the real Runner backed by os/exec.
type Exec struct {
	// Stdio, when true, attaches the child to the current terminal (for
	// interactive tools like `az login` device-code or huh-free prompts).
	Stdio bool
}

// Run executes name with args. env entries (KEY=VALUE) are appended to the
// current environment. On failure the error includes stderr.
func (e Exec) Run(ctx context.Context, env []string, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = append(os.Environ(), env...)

	if e.Stdio {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("%s: %w", name, err)
		}
		return "", nil
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%s %v: %w: %s", name, args, err, stderr.String())
	}
	return stdout.String(), nil
}
