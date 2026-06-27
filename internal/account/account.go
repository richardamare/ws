// Package account distinguishes the personal vs. work (work) identity context
// a project belongs to. Repos live under ~/Developer/Personal and
// ~/Developer/work, and work is a different Azure tenant — so the Reader-SP
// naming (and which tenant to log into) depends on the project's root.
package account

import (
	"path/filepath"
	"strings"
)

// Account is an identity context.
type Account struct {
	Name     string // "personal" | "work"
	SPPrefix string // service-principal name prefix for `ws new`
}

var (
	Personal = Account{Name: "personal", SPPrefix: "sp-ws"}
	work   = Account{Name: "work", SPPrefix: "sp-ws"}
	Unknown  = Account{Name: "unknown", SPPrefix: "sp-ws"}
)

// FromPath infers the account from a project's working directory by looking for
// a /Developer/<root>/ segment (case-insensitive). Unknown if neither matches.
func FromPath(cwd string) Account {
	p := strings.ToLower(filepath.Clean(cwd))
	switch {
	case strings.Contains(p, "/developer/work"):
		return work
	case strings.Contains(p, "/developer/personal"):
		return Personal
	default:
		return Unknown
	}
}

// SPName builds the Reader service-principal name for a project under this
// account, e.g. "sp-ws-proj1-reader".
func (a Account) SPName(project string) string {
	return a.SPPrefix + "-" + project + "-reader"
}
