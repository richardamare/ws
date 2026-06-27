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

func TestSurfaceArgsBrowser(t *testing.T) {
	args := surfaceArgs("workspace:1", &config.Project{}, config.Tab{Type: "browser", URL: "https://x"})
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--type browser") || !strings.Contains(joined, "--url https://x") {
		t.Errorf("got %q", joined)
	}
	if !strings.Contains(joined, "--workspace workspace:1") {
		t.Errorf("must target the workspace ref: %q", joined)
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
