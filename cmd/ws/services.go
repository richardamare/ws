package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/richardamare/ws/internal/azure"
	"github.com/richardamare/ws/internal/cmux"
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/run"
)

// bg is the context for external commands. Background is fine: az/cmux manage
// their own timeouts and az login device-code is intentionally long-lived.
func bg() context.Context { return context.Background() }

func itoa(n int) string { return strconv.Itoa(n) }

func azureSvc() azure.Service { return azure.Service{Run: run.Exec{}} }
func cmuxSvc() cmux.Service   { return cmux.Service{Run: run.Exec{}} }

// scopedEnv returns the AZURE_CONFIG_DIR entry for a project's isolated login,
// or nil, so setup/teardown scripts run with the same scoped Azure context as
// the workspace tabs.
func scopedEnv(p *config.Project) []string {
	if p.Azure != nil && p.Azure.ConfigDir != "" {
		return []string{"AZURE_CONFIG_DIR=" + config.ExpandHome(p.Azure.ConfigDir)}
	}
	return nil
}

// runScript executes a project's setup/teardown commands as shell commands in
// the project cwd, with the scoped env. Output streams to the terminal unless
// quiet (e.g. under --json). It does not open tabs.
func runScript(ctx context.Context, p *config.Project, cmds []string, quiet bool) error {
	if len(cmds) == 0 {
		return nil
	}
	cwd := config.ExpandHome(p.Cwd)
	r := run.Exec{Stdio: !quiet}
	for _, c := range cmds {
		line := c
		if cwd != "" {
			line = "cd '" + cwd + "' && " + c
		}
		if _, err := r.Run(ctx, scopedEnv(p), "sh", "-c", line); err != nil {
			return fmt.Errorf("script %q: %w", c, err)
		}
	}
	return nil
}

// stdio runs an interactive tool attached to the terminal (claude, az login).
func stdio(ctx context.Context, name string, args ...string) error {
	_, err := run.Exec{Stdio: true}.Run(ctx, nil, name, args...)
	return err
}

// certDir returns ~/.config/ws/certs, creating it.
func certDir() (string, error) {
	projDir, err := config.DefaultDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(filepath.Dir(projDir), "certs")
	return dir, os.MkdirAll(dir, 0o700)
}

// copyFile copies src to dst with 0600 perms.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
