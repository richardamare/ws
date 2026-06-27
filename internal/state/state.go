// Package state persists small runtime facts that are not user-edited config —
// chiefly the cmux workspace ref opened for a project, so `ws down` can close
// it. Stored as one JSON file per project under ~/.config/ws/state/. Not a
// database; safe to delete.
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type record struct {
	WorkspaceRef string `json:"workspace_ref,omitempty"`
}

// Store is a directory of per-project state files.
type Store struct {
	Dir string
}

// NewStore returns a Store under ~/.config/ws/state (honouring XDG_CONFIG_HOME).
func NewStore() (*Store, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, ".config")
	}
	return &Store{Dir: filepath.Join(dir, "ws", "state")}, nil
}

func (s *Store) path(name string) string { return filepath.Join(s.Dir, name+".json") }

func (s *Store) load(name string) record {
	var r record
	if b, err := os.ReadFile(s.path(name)); err == nil {
		_ = json.Unmarshal(b, &r)
	}
	return r
}

func (s *Store) save(name string, r record) error {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(name), b, 0o644)
}

// SetWorkspaceRef records the cmux workspace ref opened for a project.
func (s *Store) SetWorkspaceRef(name, ref string) error {
	r := s.load(name)
	r.WorkspaceRef = ref
	return s.save(name, r)
}

// WorkspaceRef returns the recorded ref, or "" if none.
func (s *Store) WorkspaceRef(name string) string {
	return s.load(name).WorkspaceRef
}

// Clear removes a project's state file.
func (s *Store) Clear(name string) error {
	err := os.Remove(s.path(name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
