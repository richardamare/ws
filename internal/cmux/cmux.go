// Package cmux drives the cmux terminal to open a project's workspace and tabs.
// ws shells out to the `cmux` binary (the socket is a fallback). Terminals that
// belong to an Azure project export the project's scoped AZURE_CONFIG_DIR so the
// tab inherits the Reader login, never the personal admin token.
package cmux

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/run"
)

// cmux prints replies like "OK workspace:13"; extract the handle (a short ref
// like workspace:13 or a UUID) from whatever it returns.
var refPattern = regexp.MustCompile(`(?:window|workspace|pane|surface):\d+|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

func parseHandle(out string) string {
	if m := refPattern.FindString(out); m != "" {
		return m
	}
	// Fallback: last whitespace-delimited token, sans a leading "OK".
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

// Service wraps the cmux CLI.
type Service struct {
	Run run.Runner
}

// workspaceArgs builds `cmux new-workspace` for a project.
func workspaceArgs(p *config.Project) []string {
	args := []string{"new-workspace", "--name", p.Name}
	if p.Cwd != "" {
		args = append(args, "--cwd", config.ExpandHome(p.Cwd))
	}
	return append(args, "--focus", "true")
}

// terminalCommand returns the shell command for a terminal tab, prefixed with
// the scoped AZURE_CONFIG_DIR export when the project has Azure config.
func terminalCommand(p *config.Project, tab config.Tab) string {
	var prefix string
	if p.Azure != nil && p.Azure.ConfigDir != "" {
		prefix = fmt.Sprintf("export AZURE_CONFIG_DIR=%q; ", config.ExpandHome(p.Azure.ConfigDir))
	}
	if tab.Run == "" {
		if prefix == "" {
			return ""
		}
		return prefix + "exec $SHELL"
	}
	return prefix + tab.Run
}

// surfaceArgs builds `cmux new-surface` for one tab in workspace ref.
func surfaceArgs(ref string, p *config.Project, tab config.Tab) []string {
	args := []string{"new-surface", "--workspace", ref}
	if tab.Type == "browser" {
		return append(args, "--type", "browser", "--url", tab.URL)
	}
	args = append(args, "--type", "terminal")
	if cmd := terminalCommand(p, tab); cmd != "" {
		args = append(args, "--command", cmd)
	}
	return args
}

// Open creates the workspace and its tabs, returning the workspace ref.
func (s Service) Open(ctx context.Context, p *config.Project) (string, error) {
	out, err := s.Run.Run(ctx, nil, "cmux", workspaceArgs(p)...)
	if err != nil {
		return "", err
	}
	ref := parseHandle(out)
	if ref == "" {
		return "", fmt.Errorf("could not parse workspace ref from cmux output: %q", out)
	}
	for _, tab := range p.Tabs {
		if _, err := s.Run.Run(ctx, nil, "cmux", surfaceArgs(ref, p, tab)...); err != nil {
			return ref, fmt.Errorf("open tab %q: %w", tab.Name, err)
		}
	}
	return ref, nil
}

// Close closes a workspace by ref.
func (s Service) Close(ctx context.Context, ref string) error {
	if ref == "" {
		return fmt.Errorf("no workspace ref recorded")
	}
	_, err := s.Run.Run(ctx, nil, "cmux", "close-workspace", "--workspace", ref)
	return err
}

// NewTerminal opens a terminal surface running command in the current workspace.
// Used by `ws elevate` to spawn a marked personal-admin tab.
func (s Service) NewTerminal(ctx context.Context, command string) error {
	args := []string{"new-surface", "--type", "terminal", "--focus", "true"}
	if command != "" {
		args = append(args, "--command", command)
	}
	_, err := s.Run.Run(ctx, nil, "cmux", args...)
	return err
}

// ReloadConfig reloads cmux.json (and Ghostty config) in place — no app restart.
func (s Service) ReloadConfig(ctx context.Context) error {
	_, err := s.Run.Run(ctx, nil, "cmux", "reload-config")
	return err
}

// ValidateConfig checks cmux.json is syntactically valid.
func (s Service) ValidateConfig(ctx context.Context) error {
	_, err := s.Run.Run(ctx, nil, "cmux", "config", "validate")
	return err
}

// ResumeID returns the focused agent surface's resume id, for session bookmarks.
func (s Service) ResumeID(ctx context.Context) (string, error) {
	out, err := s.Run.Run(ctx, nil, "cmux", "surface", "resume", "show")
	if err != nil {
		return "", err
	}
	s2 := strings.TrimSpace(out)
	s2 = strings.TrimPrefix(s2, "OK ")
	return strings.TrimSpace(s2), nil
}
