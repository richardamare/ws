package cmux

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"

	"github.com/richardamare/ws/internal/config"
)

// A cmux.json `commands[]` entry: a reusable workspace template. ws generates
// one per project so cmux's own restore reopens the tabs after a crash/close,
// without ws running (the durable half of the hybrid model, ADR-0003).

type Surface struct {
	Type    string `json:"type"`
	Name    string `json:"name,omitempty"`
	Command string `json:"command,omitempty"`
	URL     string `json:"url,omitempty"`
	Focus   bool   `json:"focus,omitempty"`
}

type pane struct {
	Surfaces []Surface `json:"surfaces"`
}

type child struct {
	Pane pane `json:"pane"`
}

type layout struct {
	Direction string  `json:"direction"`
	Children  []child `json:"children"`
}

type workspace struct {
	Name   string `json:"name"`
	Cwd    string `json:"cwd,omitempty"`
	Color  string `json:"color,omitempty"`
	Layout layout `json:"layout"`
}

// Command is one cmux.json commands[] entry.
type Command struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Keywords    []string  `json:"keywords,omitempty"`
	Workspace   workspace `json:"workspace"`
}

var palette = []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

func colorFor(name string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return palette[int(h.Sum32())%len(palette)]
}

// BuildCommand turns a project into a cmux workspace-template command. One pane
// per tab, laid out horizontally. Terminal panes inherit the scoped
// AZURE_CONFIG_DIR via terminalCommand.
func BuildCommand(p *config.Project) Command {
	children := make([]child, 0, len(p.Tabs))
	for i, tab := range p.Tabs {
		s := Surface{Type: tab.Type, Name: tab.Name, Focus: i == 0}
		if tab.Type == "browser" {
			s.URL = tab.URL
		} else {
			s.Type = "terminal"
			s.Command = terminalCommand(p, tab)
		}
		children = append(children, child{Pane: pane{Surfaces: []Surface{s}}})
	}
	cwd := ""
	if p.Cwd != "" {
		cwd = config.ExpandHome(p.Cwd)
	}
	return Command{
		Name:        p.Name,
		Description: "ws workspace for " + p.Name,
		Keywords:    []string{"ws", p.Name},
		Workspace: workspace{
			Name:   p.Name,
			Cwd:    cwd,
			Color:  colorFor(p.Name),
			Layout: layout{Direction: "horizontal", Children: children},
		},
	}
}

// DefaultConfigPath returns ~/.config/cmux/cmux.json.
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "cmux", "cmux.json"), nil
}

// UpsertCommand merges c into the commands[] array of the cmux.json at path,
// matching by name. The existing file is backed up to <path>.bak first. A
// missing file is created with a minimal valid config.
func UpsertCommand(path string, c Command) error {
	root := map[string]any{
		"$schema":       "https://raw.githubusercontent.com/manaflow-ai/cmux/main/web/data/cmux.schema.json",
		"schemaVersion": 1,
	}

	if b, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(b, &root); err != nil {
			return fmt.Errorf("cmux.json is not valid JSON (%w); fix or remove it so ws can manage it", err)
		}
		_ = os.WriteFile(path+".bak", b, 0o644) // best-effort backup
	}

	// Decode existing commands, replace-or-append by name.
	var commands []Command
	if raw, ok := root["commands"]; ok {
		if b, err := json.Marshal(raw); err == nil {
			_ = json.Unmarshal(b, &commands)
		}
	}
	replaced := false
	for i := range commands {
		if commands[i].Name == c.Name {
			commands[i] = c
			replaced = true
			break
		}
	}
	if !replaced {
		commands = append(commands, c)
	}
	root["commands"] = commands

	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o644)
}
