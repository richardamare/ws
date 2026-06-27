package cmux

import (
	"strings"
	"testing"

	"github.com/richardamare/ws/internal/config"
)

func TestTerminalCommandScopedEnv(t *testing.T) {
	p := &config.Project{Azure: &config.Azure{ConfigDir: "~/.azure-proj1"}}

	cmd := terminalCommand(p, config.Tab{Type: "terminal", Run: "claude"})
	if !strings.Contains(cmd, "AZURE_CONFIG_DIR=") || !strings.HasSuffix(cmd, "claude") {
		t.Errorf("expected scoped env + claude, got %q", cmd)
	}

	shell := terminalCommand(p, config.Tab{Type: "terminal"})
	if !strings.Contains(shell, "exec $SHELL") {
		t.Errorf("empty-run terminal should launch a shell, got %q", shell)
	}
}

func TestTerminalCommandNoAzure(t *testing.T) {
	p := &config.Project{}
	if got := terminalCommand(p, config.Tab{Type: "terminal"}); got != "" {
		t.Errorf("no azure + no run => empty command, got %q", got)
	}
	if got := terminalCommand(p, config.Tab{Type: "terminal", Run: "vim"}); got != "vim" {
		t.Errorf("got %q", got)
	}
}

func TestNewSurfaceArgsBrowser(t *testing.T) {
	args := newSurfaceArgs("workspace:1", &config.Project{}, config.Tab{Type: "browser", URL: "https://x"}, true)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--type browser") || !strings.Contains(joined, "--url https://x") {
		t.Errorf("got %q", joined)
	}
	if !strings.Contains(joined, "--workspace workspace:1") || !strings.Contains(joined, "--focus true") {
		t.Errorf("must target the workspace ref and focus: %q", joined)
	}
}

func TestNewSurfaceArgsTerminalCwd(t *testing.T) {
	p := &config.Project{Cwd: "/tmp/proj"}
	args := newSurfaceArgs("workspace:2", p, config.Tab{Type: "terminal", Name: "Shell"}, false)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--type terminal") || !strings.Contains(joined, "--working-directory /tmp/proj") {
		t.Errorf("terminal should set working dir: %q", joined)
	}
	if strings.Contains(joined, "--command") {
		t.Errorf("new-surface has no --command flag: %q", joined)
	}
}

func TestWorkspaceArgsEnv(t *testing.T) {
	p := &config.Project{Name: "p", Azure: &config.Azure{ConfigDir: "~/.azure-p"}}
	joined := strings.Join(workspaceArgs(p), " ")
	if !strings.Contains(joined, "--env AZURE_CONFIG_DIR=") {
		t.Errorf("workspace should pass scoped env: %q", joined)
	}
}

func TestParseHandle(t *testing.T) {
	cases := map[string]string{
		"OK workspace:13": "workspace:13",
		"workspace:1":     "workspace:1",
		"OK 3ee36000-a8a6-4a5a-ae8c-3d363953d18e":  "3ee36000-a8a6-4a5a-ae8c-3d363953d18e",
		"created surface surface:7 in workspace:2": "surface:7",
		"plainvalue": "plainvalue",
	}
	for in, want := range cases {
		if got := parseHandle(in); got != want {
			t.Errorf("parseHandle(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestWorkspaceArgs(t *testing.T) {
	args := workspaceArgs(&config.Project{Name: "proj1", Cwd: "/tmp/x"})
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "new-workspace") || !strings.Contains(joined, "--name proj1") {
		t.Errorf("got %q", joined)
	}
}
