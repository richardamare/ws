package config

import (
	"path/filepath"
	"testing"
)

func newTempStore(t *testing.T) *Store {
	t.Helper()
	return &Store{Dir: filepath.Join(t.TempDir(), "projects")}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	s := newTempStore(t)
	want := &Project{
		Name: "proj1",
		Cwd:  "~/code/proj1",
		Azure: &Azure{
			SPAppID:       "app-123",
			Tenant:        "tenant-1",
			Cert:          "~/.config/ws/certs/proj1.pem",
			ConfigDir:     "~/.azure-proj1",
			Subscription:  "sub-1",
			ResourceGroup: "rg-proj1",
		},
		Tabs: []Tab{
			{Type: "terminal", Name: "Claude", Run: "claude"},
			{Type: "browser", Name: "Repo", URL: "https://github.com/me/proj1"},
		},
		Sessions: []Bookmark{{Label: "auth", ID: "3ee3", Note: "rbac"}},
		Setup:    []string{"docker compose up -d"},
		Teardown: []string{"docker compose down"},
	}
	if err := s.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := s.Load("proj1")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Name != want.Name || got.Cwd != want.Cwd {
		t.Errorf("scalar mismatch: got %+v", got)
	}
	if got.Azure == nil || got.Azure.ResourceGroup != "rg-proj1" {
		t.Errorf("azure not round-tripped: %+v", got.Azure)
	}
	if len(got.Tabs) != 2 || got.Tabs[1].URL != "https://github.com/me/proj1" {
		t.Errorf("tabs not round-tripped: %+v", got.Tabs)
	}
	if len(got.Sessions) != 1 || got.Sessions[0].Label != "auth" {
		t.Errorf("sessions not round-tripped: %+v", got.Sessions)
	}
	if len(got.Setup) != 1 || got.Setup[0] != "docker compose up -d" {
		t.Errorf("setup not round-tripped: %+v", got.Setup)
	}
	if len(got.Teardown) != 1 || got.Teardown[0] != "docker compose down" {
		t.Errorf("teardown not round-tripped: %+v", got.Teardown)
	}
}

func TestSaveRequiresName(t *testing.T) {
	s := newTempStore(t)
	if err := s.Save(&Project{}); err == nil {
		t.Fatal("expected error saving project with empty name")
	}
}

func TestListSortedAndEmpty(t *testing.T) {
	s := newTempStore(t)
	names, err := s.List()
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected no projects, got %v", names)
	}
	for _, n := range []string{"beta", "alpha"} {
		if err := s.Save(&Project{Name: n}); err != nil {
			t.Fatal(err)
		}
	}
	names, err = s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("expected [alpha beta], got %v", names)
	}
}

func TestLoadMissing(t *testing.T) {
	s := newTempStore(t)
	if _, err := s.Load("nope"); err == nil {
		t.Fatal("expected error loading missing project")
	}
}
