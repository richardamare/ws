// Package config loads and saves per-project ws configuration. One YAML file
// per project under ~/.config/ws/projects/. There is no database.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Project is the on-disk shape of a single project's config.
type Project struct {
	Name string `yaml:"name"`
	Cwd  string `yaml:"cwd,omitempty"`

	Azure     *Azure     `yaml:"azure,omitempty"`
	Tabs      []Tab      `yaml:"tabs,omitempty"`
	Sessions  []Bookmark `yaml:"sessions,omitempty"`
	Container *Container `yaml:"container,omitempty"`
}

// Azure holds the scoped Reader service-principal login for a project. The SP is
// Reader-only on a single resource group; see docs/security/README.md.
type Azure struct {
	SPAppID       string `yaml:"sp_app_id"`
	Tenant        string `yaml:"tenant"`
	Cert          string `yaml:"cert"`
	ConfigDir     string `yaml:"config_dir"`
	Subscription  string `yaml:"subscription"`
	ResourceGroup string `yaml:"resource_group"`
}

// Tab is one cmux surface to open: a terminal (optionally running Run) or a
// browser pointed at URL.
type Tab struct {
	Type string `yaml:"type"` // "terminal" | "browser"
	Name string `yaml:"name"`
	Run  string `yaml:"run,omitempty"`
	URL  string `yaml:"url,omitempty"`
}

// Bookmark is a curated Claude Code session, reused via `ws resume`.
type Bookmark struct {
	Label string `yaml:"label"`
	ID    string `yaml:"id"`
	Note  string `yaml:"note,omitempty"`
}

// Container is the optional dev-container block. Absent means host-only.
type Container struct {
	Compose   string `yaml:"compose"`
	Service   string `yaml:"service"`
	ExecShell string `yaml:"exec_shell,omitempty"`
}

// ExpandHome replaces a leading ~ with the user's home directory. Paths without
// a leading ~ are returned unchanged.
func ExpandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}

// Store is a directory of project YAML files.
type Store struct {
	Dir string
}

// DefaultDir returns ~/.config/ws/projects (honouring XDG_CONFIG_HOME).
func DefaultDir() (string, error) {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "ws", "projects"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ws", "projects"), nil
}

// NewStore returns a Store rooted at the default project directory.
func NewStore() (*Store, error) {
	dir, err := DefaultDir()
	if err != nil {
		return nil, err
	}
	return &Store{Dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.Dir, name+".yaml")
}

// Load reads one project by name.
func (s *Store) Load(name string) (*Project, error) {
	b, err := os.ReadFile(s.path(name))
	if err != nil {
		return nil, err
	}
	var p Project
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	return &p, nil
}

// Save writes a project, creating the store directory if needed.
func (s *Store) Save(p *Project) error {
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(p.Name), b, 0o644)
}

// Remove deletes a project's config file. Missing is not an error.
func (s *Store) Remove(name string) error {
	err := os.Remove(s.path(name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// List returns the project names found in the store, sorted.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if n := strings.TrimSuffix(e.Name(), ".yaml"); n != e.Name() {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names, nil
}
