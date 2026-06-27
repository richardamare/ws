// Package azure drives the Azure CLI for ws. Everyday sessions log in as a
// per-project Reader service principal, isolated in the project's own
// AZURE_CONFIG_DIR so they can never see the personal admin token. The SP is
// Reader-only on a single resource group. Write/Terraform is never done here —
// that is the deliberate `ws elevate` path. See docs/security/README.md.
package azure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/run"
)

// Service wraps the az CLI.
type Service struct {
	Run run.Runner
}

// Identity is the subset of `az account show` we care about.
type Identity struct {
	User         string // the SP appId (or user UPN)
	TenantID     string
	Subscription string
}

// env returns the AZURE_CONFIG_DIR isolation for a project's scoped login.
func env(a *config.Azure) []string {
	return []string{"AZURE_CONFIG_DIR=" + config.ExpandHome(a.ConfigDir)}
}

// loginArgs builds the non-interactive service-principal login (cert-based).
func loginArgs(a *config.Azure) []string {
	return []string{
		"login", "--service-principal",
		"-u", a.SPAppID,
		"--certificate", config.ExpandHome(a.Cert),
		"--tenant", a.Tenant,
		"--allow-no-subscriptions",
	}
}

// createReaderSPArgs builds the `ws new` service-principal creation, scoped to a
// single resource group with the Reader role and a generated certificate.
func createReaderSPArgs(name, subscription, rg string) []string {
	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscription, rg)
	return []string{
		"ad", "sp", "create-for-rbac",
		"--name", name,
		"--role", "Reader",
		"--scopes", scope,
		"--create-cert",
	}
}

// Status returns the identity currently logged into the project's config dir, or
// an error if none.
func (s Service) Status(ctx context.Context, a *config.Azure) (Identity, error) {
	out, err := s.Run.Run(ctx, env(a), "az", "account", "show", "-o", "json")
	if err != nil {
		return Identity{}, err
	}
	var raw struct {
		TenantID string `json:"tenantId"`
		ID       string `json:"id"`
		User     struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"user"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return Identity{}, fmt.Errorf("parse az account show: %w", err)
	}
	return Identity{
		User:         raw.User.Name,
		TenantID:     raw.TenantID,
		Subscription: raw.ID,
	}, nil
}

// Login ensures the project's Reader SP is logged into its config dir. Idempotent:
// if the expected SP is already active, it does nothing.
func (s Service) Login(ctx context.Context, a *config.Azure) error {
	if a == nil {
		return fmt.Errorf("project has no azure config")
	}
	if id, err := s.Status(ctx, a); err == nil && id.User == a.SPAppID {
		return nil // already logged in as the right SP
	}
	_, err := s.Run.Run(ctx, env(a), "az", loginArgs(a)...)
	return err
}

// SPResult is the parsed output of an SP create / credential reset.
type SPResult struct {
	AppID    string `json:"appId"`
	Tenant   string `json:"tenant"`
	CertFile string `json:"fileWithCertAndPrivateKey"`
}

// ParseSPCreate parses `az ad sp create-for-rbac`/`credential reset` JSON.
func ParseSPCreate(jsonOut string) (SPResult, error) {
	var r SPResult
	if err := json.Unmarshal([]byte(jsonOut), &r); err != nil {
		return SPResult{}, fmt.Errorf("parse sp result: %w", err)
	}
	if r.AppID == "" || r.CertFile == "" {
		return SPResult{}, fmt.Errorf("sp result missing appId or cert file")
	}
	return r, nil
}

// CreateReaderSP creates the scoped Reader service principal for a new project.
func (s Service) CreateReaderSP(ctx context.Context, name, subscription, rg string) (SPResult, error) {
	out, err := s.Run.Run(ctx, nil, "az", createReaderSPArgs(name, subscription, rg)...)
	if err != nil {
		return SPResult{}, err
	}
	return ParseSPCreate(out)
}

// RotateCert issues a fresh certificate for an existing SP.
func (s Service) RotateCert(ctx context.Context, appID string) (SPResult, error) {
	out, err := s.Run.Run(ctx, nil, "az",
		"ad", "sp", "credential", "reset", "--id", appID, "--create-cert")
	if err != nil {
		return SPResult{}, err
	}
	return ParseSPCreate(out)
}
