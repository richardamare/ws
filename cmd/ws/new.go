package main

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/richardamare/ws/internal/account"
	"github.com/richardamare/ws/internal/config"
	"github.com/richardamare/ws/internal/output"
	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	var cwd, sub, rg, repo string
	cmd := &cobra.Command{
		Use:   "new [project]",
		Short: "Create a project: scoped Reader SP + cert + config",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}

			if interactive() {
				if err := promptMissing(&name, &cwd, &sub, &rg, &repo); err != nil {
					return err
				}
			}
			if name == "" || cwd == "" || sub == "" || rg == "" {
				return fmt.Errorf("name, --cwd, --sub and --rg are required")
			}

			store, err := config.NewStore()
			if err != nil {
				return err
			}
			if _, err := store.Load(name); err == nil {
				return fmt.Errorf("project %q already exists", name)
			}

			acct := account.FromPath(cwd)
			spName := acct.SPName(name)
			sp, err := azureSvc().CreateReaderSP(bg(), spName, sub, rg)
			if err != nil {
				return fmt.Errorf("create reader SP: %w", err)
			}

			dir, err := certDir()
			if err != nil {
				return err
			}
			certPath := filepath.Join(dir, name+".pem")
			if err := copyFile(sp.CertFile, certPath); err != nil {
				return fmt.Errorf("store cert: %w", err)
			}

			p := &config.Project{
				Name: name,
				Cwd:  cwd,
				Azure: &config.Azure{
					SPAppID:       sp.AppID,
					Tenant:        sp.Tenant,
					Cert:          certPath,
					ConfigDir:     "~/.azure-" + name,
					Subscription:  sub,
					ResourceGroup: rg,
				},
				Tabs: defaultTabs(repo),
			}
			if err := store.Save(p); err != nil {
				return err
			}

			return output.Record(cmd.OutOrStdout(), resolveFormat(), []output.KV{
				{Key: "project", Value: name},
				{Key: "account", Value: acct.Name},
				{Key: "sp", Value: spName},
				{Key: "rg", Value: rg},
				{Key: "cert", Value: certPath},
			})
		},
	}
	cmd.Flags().StringVar(&cwd, "cwd", "", "project working directory")
	cmd.Flags().StringVar(&sub, "sub", "", "Azure subscription id")
	cmd.Flags().StringVar(&rg, "rg", "", "resource group to scope the Reader SP to")
	cmd.Flags().StringVar(&repo, "repo", "", "GitHub repo URL for a browser tab")
	return cmd
}

func defaultTabs(repo string) []config.Tab {
	tabs := []config.Tab{
		{Type: "terminal", Name: "Claude", Run: "claude"},
		{Type: "terminal", Name: "Shell"},
	}
	if repo != "" {
		tabs = append(tabs, config.Tab{Type: "browser", Name: "Repo", URL: repo})
	}
	return tabs
}

// promptMissing fills empty fields via huh inputs (interactive only).
func promptMissing(name, cwd, sub, rg, repo *string) error {
	fields := []huh.Field{}
	add := func(title string, v *string, required bool) {
		if *v != "" {
			return
		}
		in := huh.NewInput().Title(title).Value(v)
		if required {
			in = in.Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("required")
				}
				return nil
			})
		}
		fields = append(fields, in)
	}
	add("Project name", name, true)
	add("Working directory", cwd, true)
	add("Azure subscription id", sub, true)
	add("Resource group", rg, true)
	add("GitHub repo URL (optional)", repo, false)
	if len(fields) == 0 {
		return nil
	}
	return huh.NewForm(huh.NewGroup(fields...)).Run()
}
