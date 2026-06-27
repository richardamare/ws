// Package workspace turns a project config into the ordered set of actions
// `ws up` will perform (scoped Azure login, optional container, cmux tabs). The
// plan is pure data so it can be rendered, dry-run, and tested without touching
// cmux or Azure. Execution of the plan lands in follow-up work.
package workspace

import (
	"fmt"

	"github.com/richardamare/ws/internal/config"
)

// Step is one action in an `up` plan.
type Step struct {
	Action string
	Detail string
}

// Plan is the ordered list of steps for a project.
type Plan struct {
	Project string
	Steps   []Step
}

// Plan computes the steps for a project, in execution order.
func PlanFor(p *config.Project) Plan {
	plan := Plan{Project: p.Name}

	if p.Azure != nil {
		plan.Steps = append(plan.Steps, Step{
			Action: "az-login",
			Detail: fmt.Sprintf("Reader SP %s into %s (Reader on %s)",
				p.Azure.SPAppID, p.Azure.ConfigDir, p.Azure.ResourceGroup),
		})
	}

	if p.Container != nil {
		plan.Steps = append(plan.Steps, Step{
			Action: "container",
			Detail: fmt.Sprintf("docker compose -f %s up -d (service %s)",
				p.Container.Compose, p.Container.Service),
		})
	}

	plan.Steps = append(plan.Steps, Step{
		Action: "cmux-workspace",
		Detail: fmt.Sprintf("open workspace %q with %d tab(s)", p.Name, len(p.Tabs)),
	})
	for _, t := range p.Tabs {
		detail := t.Name
		switch t.Type {
		case "terminal":
			if t.Run != "" {
				detail = fmt.Sprintf("terminal %q runs %q", t.Name, t.Run)
			} else {
				detail = fmt.Sprintf("terminal %q", t.Name)
			}
		case "browser":
			detail = fmt.Sprintf("browser %q -> %s", t.Name, t.URL)
		}
		plan.Steps = append(plan.Steps, Step{Action: "tab", Detail: detail})
	}

	return plan
}

// Rows renders the plan as table rows for internal/output.
func (p Plan) Rows() [][]string {
	rows := make([][]string, 0, len(p.Steps))
	for _, s := range p.Steps {
		rows = append(rows, []string{s.Action, s.Detail})
	}
	return rows
}
