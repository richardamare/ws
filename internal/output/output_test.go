package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	cases := []struct {
		name            string
		jsonFlag, plain bool
		isTTY           bool
		want            Format
	}{
		{"json wins over everything", true, true, true, JSON},
		{"plain flag on tty", false, true, true, Plain},
		{"non-tty defaults to plain", false, false, false, Plain},
		{"tty defaults to pretty", false, false, true, Pretty},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Detect(c.jsonFlag, c.plain, c.isTTY); got != c.want {
				t.Errorf("Detect = %v, want %v", got, c.want)
			}
		})
	}
}

func TestRecordPlain(t *testing.T) {
	var b bytes.Buffer
	fields := []KV{{"name", "proj1"}, {"status", "up"}}
	if err := Record(&b, Plain, fields); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, "name:") || !strings.Contains(out, "proj1") {
		t.Errorf("plain record missing data:\n%s", out)
	}
}

func TestRecordJSON(t *testing.T) {
	var b bytes.Buffer
	if err := Record(&b, JSON, []KV{{"name", "proj1"}}); err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal(b.Bytes(), &m); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, b.String())
	}
	if m["name"] != "proj1" {
		t.Errorf("got %v", m)
	}
}

func TestTableJSON(t *testing.T) {
	var b bytes.Buffer
	headers := []string{"name", "status"}
	rows := [][]string{{"proj1", "up"}, {"proj2", "down"}}
	if err := Table(&b, JSON, headers, rows); err != nil {
		t.Fatal(err)
	}
	var objs []map[string]string
	if err := json.Unmarshal(b.Bytes(), &objs); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(objs) != 2 || objs[0]["name"] != "proj1" || objs[1]["status"] != "down" {
		t.Errorf("got %v", objs)
	}
}

func TestTablePlainAligned(t *testing.T) {
	var b bytes.Buffer
	if err := Table(&b, Plain, []string{"LABEL", "ID"}, [][]string{{"auth", "3ee3"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "LABEL") || !strings.Contains(b.String(), "auth") {
		t.Errorf("missing content:\n%s", b.String())
	}
}
