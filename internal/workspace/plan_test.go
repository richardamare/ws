package workspace

import (
	"strings"
	"testing"

	"github.com/richardamare/ws/internal/config"
)

func TestPlanForOrderAndContent(t *testing.T) {
	p := &config.Project{
		Name: "proj1",
		Azure: &config.Azure{
			SPAppID:       "app-123",
			ConfigDir:     "~/.azure-proj1",
			ResourceGroup: "rg-proj1",
		},
		Container: &config.Container{Compose: "docker-compose.yml", Service: "devcontainer"},
		Tabs: []config.Tab{
			{Type: "terminal", Name: "Claude", Run: "claude"},
			{Type: "browser", Name: "Repo", URL: "https://github.com/me/proj1"},
		},
	}
	plan := PlanFor(p)

	if plan.Steps[0].Action != "az-login" {
		t.Errorf("expected az-login first, got %q", plan.Steps[0].Action)
	}
	if plan.Steps[1].Action != "container" {
		t.Errorf("expected container second, got %q", plan.Steps[1].Action)
	}
	if plan.Steps[2].Action != "cmux-workspace" {
		t.Errorf("expected cmux-workspace third, got %q", plan.Steps[2].Action)
	}
	if !strings.Contains(plan.Steps[0].Detail, "rg-proj1") {
		t.Errorf("az-login detail should mention the RG: %q", plan.Steps[0].Detail)
	}
	// one row per tab after the workspace step
	tabRows := 0
	for _, s := range plan.Steps {
		if s.Action == "tab" {
			tabRows++
		}
	}
	if tabRows != 2 {
		t.Errorf("expected 2 tab steps, got %d", tabRows)
	}
}

func TestPlanForMinimal(t *testing.T) {
	plan := PlanFor(&config.Project{Name: "bare"})
	if len(plan.Steps) != 1 || plan.Steps[0].Action != "cmux-workspace" {
		t.Fatalf("minimal project should yield only the workspace step, got %+v", plan.Steps)
	}
	if len(plan.Rows()) != 1 {
		t.Errorf("Rows length mismatch: %v", plan.Rows())
	}
}
