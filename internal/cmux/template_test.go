package cmux

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardamare/ws/internal/config"
)

func sampleProject() *config.Project {
	return &config.Project{
		Name:  "proj1",
		Cwd:   "/tmp/proj1",
		Azure: &config.Azure{ConfigDir: "~/.azure-proj1"},
		Tabs: []config.Tab{
			{Type: "terminal", Name: "Claude", Run: "claude"},
			{Type: "browser", Name: "Repo", URL: "https://github.com/me/proj1"},
		},
	}
}

func TestBuildCommandShape(t *testing.T) {
	c := BuildCommand(sampleProject())
	if c.Name != "proj1" || c.Workspace.Name != "proj1" {
		t.Fatalf("names wrong: %+v", c)
	}
	if len(c.Workspace.Layout.Children) != 2 {
		t.Fatalf("expected 2 panes, got %d", len(c.Workspace.Layout.Children))
	}
	term := c.Workspace.Layout.Children[0].Pane.Surfaces[0]
	if term.Type != "terminal" || !term.Focus {
		t.Errorf("first surface should be focused terminal: %+v", term)
	}
	if term.Command == "" || term.Command[:6] != "export" {
		t.Errorf("terminal should export scoped env: %q", term.Command)
	}
	browser := c.Workspace.Layout.Children[1].Pane.Surfaces[0]
	if browser.Type != "browser" || browser.URL == "" {
		t.Errorf("second surface should be a browser: %+v", browser)
	}
}

func TestUpsertCommandCreatesAndReplaces(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cmux.json")

	if err := UpsertCommand(path, BuildCommand(sampleProject())); err != nil {
		t.Fatal(err)
	}
	root := readRoot(t, path)
	if root["schemaVersion"] == nil {
		t.Error("expected schemaVersion to be seeded")
	}
	if got := len(root["commands"].([]any)); got != 1 {
		t.Fatalf("expected 1 command, got %d", got)
	}

	// upsert same name → still 1, replaced (not appended)
	p2 := sampleProject()
	p2.Tabs = p2.Tabs[:1]
	if err := UpsertCommand(path, BuildCommand(p2)); err != nil {
		t.Fatal(err)
	}
	root = readRoot(t, path)
	cmds := root["commands"].([]any)
	if len(cmds) != 1 {
		t.Fatalf("upsert should replace, got %d commands", len(cmds))
	}

	// a different project appends
	p3 := sampleProject()
	p3.Name = "proj2"
	if err := UpsertCommand(path, BuildCommand(p3)); err != nil {
		t.Fatal(err)
	}
	if got := len(readRoot(t, path)["commands"].([]any)); got != 2 {
		t.Fatalf("expected 2 commands, got %d", got)
	}
	if _, err := os.Stat(path + ".bak"); err != nil {
		t.Errorf("expected a .bak backup: %v", err)
	}
}

func TestUpsertRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cmux.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := UpsertCommand(path, BuildCommand(sampleProject())); err == nil {
		t.Fatal("expected error on invalid existing cmux.json")
	}
}

func readRoot(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	return m
}
