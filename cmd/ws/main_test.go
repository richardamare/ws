package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCommandJSON(t *testing.T) {
	t.Cleanup(func() { flagJSON = false })

	root := newRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--json", "version"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(out.Bytes(), &m); err != nil {
		t.Fatalf("version output is not valid JSON: %v\n%s", err, out.String())
	}
	if m["version"] == "" {
		t.Errorf("expected a version field, got %v", m)
	}
}
