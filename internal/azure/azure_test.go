package azure

import (
	"context"
	"strings"
	"testing"

	"github.com/richardamare/ws/internal/config"
)

// fakeRunner records calls and returns scripted output.
type fakeRunner struct {
	calls   [][]string
	outputs []string
	errs    []error
	i       int
}

func (f *fakeRunner) Run(_ context.Context, _ []string, name string, args ...string) (string, error) {
	f.calls = append(f.calls, append([]string{name}, args...))
	var out string
	var err error
	if f.i < len(f.outputs) {
		out = f.outputs[f.i]
	}
	if f.i < len(f.errs) {
		err = f.errs[f.i]
	}
	f.i++
	return out, err
}

func sampleAzure() *config.Azure {
	return &config.Azure{
		SPAppID:       "app-123",
		Tenant:        "tenant-1",
		Cert:          "~/.config/ws/certs/proj1.pem",
		ConfigDir:     "~/.azure-proj1",
		Subscription:  "sub-1",
		ResourceGroup: "rg-proj1",
	}
}

func TestLoginArgs(t *testing.T) {
	args := loginArgs(sampleAzure())
	joined := strings.Join(args, " ")
	for _, want := range []string{"login", "--service-principal", "-u app-123", "--certificate", "--tenant tenant-1"} {
		if !strings.Contains(joined, want) {
			t.Errorf("loginArgs missing %q in %q", want, joined)
		}
	}
	if strings.Contains(joined, "~") {
		t.Errorf("cert path should be expanded, got %q", joined)
	}
}

func TestCreateReaderSPArgsScope(t *testing.T) {
	args := createReaderSPArgs("sp-x", "sub-1", "rg-1")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--role Reader") {
		t.Errorf("must be Reader role: %q", joined)
	}
	if !strings.Contains(joined, "/subscriptions/sub-1/resourceGroups/rg-1") {
		t.Errorf("scope must be the single RG: %q", joined)
	}
}

func TestLoginSkipsWhenAlreadyActive(t *testing.T) {
	// account show returns the expected SP -> no login call.
	f := &fakeRunner{outputs: []string{`{"user":{"name":"app-123","type":"servicePrincipal"},"tenantId":"tenant-1","id":"sub-1"}`}}
	if err := (Service{Run: f}).Login(context.Background(), sampleAzure()); err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 1 || f.calls[0][1] != "account" {
		t.Fatalf("expected only `az account show`, got %v", f.calls)
	}
}

func TestLoginRunsWhenNotActive(t *testing.T) {
	// account show fails (not logged in) -> a login call follows.
	f := &fakeRunner{
		outputs: []string{"", ""},
		errs:    []error{context.DeadlineExceeded, nil},
	}
	if err := (Service{Run: f}).Login(context.Background(), sampleAzure()); err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 2 {
		t.Fatalf("expected account show + login, got %v", f.calls)
	}
	if f.calls[1][1] != "login" {
		t.Errorf("second call should be login, got %v", f.calls[1])
	}
}
