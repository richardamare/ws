#!/usr/bin/env bash
# Audit hook: append one JSONL line per tool call to .claude/audit/<session>.jsonl.
# Fed Claude Code PreToolUse JSON on stdin. MUST never block a tool — always exit 0.
input=$(cat)
root="${CLAUDE_PROJECT_DIR:-$(pwd)}"
dir="$root/.claude/audit"
mkdir -p "$dir" 2>/dev/null

sid=$(printf '%s' "$input" | jq -r '.session_id // "unknown"' 2>/dev/null)
[ -z "$sid" ] && sid="unknown"
ts=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# summary = the most informative single field per tool, coerced to string, capped at 300 chars.
printf '%s' "$input" | jq -c --arg ts "$ts" '{
  ts: $ts,
  session_id: .session_id,
  cwd: .cwd,
  tool: .tool_name,
  summary: (
    (.tool_input.command
      // .tool_input.file_path
      // .tool_input.path
      // .tool_input.pattern
      // .tool_input.url
      // .tool_input.prompt
      // .tool_input.description
      // (.tool_input | keys | join(",")))
    | tostring | .[0:300]
  )
}' >> "$dir/$sid.jsonl" 2>/dev/null

exit 0
