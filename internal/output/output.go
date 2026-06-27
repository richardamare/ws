// Package output renders command results in one of three modes so ws is usable
// by both a human and an LLM: Pretty (TTY), Plain (structured text, the agent
// default on non-TTY), and JSON (strict). See docs/product/output-modes.md.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Format selects how results are rendered.
type Format int

const (
	Pretty Format = iota // human, TTY: tables + colour (colour added later)
	Plain                // structured text: labelled blocks / TSV
	JSON                 // strict machine output
)

// Detect resolves the output format from the flags and TTY state. --json wins;
// then --plain; otherwise Plain when not a TTY (best for piping into an agent),
// Pretty on a TTY.
func Detect(jsonFlag, plainFlag, isTTY bool) Format {
	switch {
	case jsonFlag:
		return JSON
	case plainFlag:
		return Plain
	case !isTTY:
		return Plain
	default:
		return Pretty
	}
}

// KV is one ordered field of a single-record result.
type KV struct {
	Key   string
	Value string
}

// Record renders a single object as key/value lines (Plain/Pretty) or a JSON
// object (JSON), preserving field order.
func Record(w io.Writer, f Format, fields []KV) error {
	if f == JSON {
		m := make(map[string]string, len(fields))
		for _, kv := range fields {
			m[kv.Key] = kv.Value
		}
		return writeJSON(w, m)
	}
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	for _, kv := range fields {
		fmt.Fprintf(tw, "%s:\t%s\n", kv.Key, kv.Value)
	}
	return tw.Flush()
}

// Table renders rows as aligned columns (Plain/Pretty) or an array of objects
// keyed by header (JSON).
func Table(w io.Writer, f Format, headers []string, rows [][]string) error {
	if f == JSON {
		objs := make([]map[string]string, 0, len(rows))
		for _, row := range rows {
			m := make(map[string]string, len(headers))
			for i, h := range headers {
				if i < len(row) {
					m[h] = row[i]
				}
			}
			objs = append(objs, m)
		}
		return writeJSON(w, objs)
	}
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	writeRow(tw, headers)
	for _, row := range rows {
		writeRow(tw, row)
	}
	return tw.Flush()
}

func writeRow(w io.Writer, cols []string) {
	for i, c := range cols {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, c)
	}
	fmt.Fprintln(w)
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
