package account

import "testing"

func TestFromPath(t *testing.T) {
	cases := map[string]Account{
		"/Users/r/Developer/Personal/ws":   Personal,
		"/Users/r/developer/personal/proj": Personal,
		"/Users/r/Developer/work/myproj": work,
		"~/Developer/work/x":             work,
		"/Users/r/code/other":              Unknown,
	}
	for path, want := range cases {
		if got := FromPath(path); got.Name != want.Name {
			t.Errorf("FromPath(%q) = %q, want %q", path, got.Name, want.Name)
		}
	}
}

func TestSPName(t *testing.T) {
	if got := Personal.SPName("proj1"); got != "sp-ws-proj1-reader" {
		t.Errorf("got %q", got)
	}
	if got := work.SPName("myproj"); got != "sp-ws-myproj-reader" {
		t.Errorf("got %q", got)
	}
}
