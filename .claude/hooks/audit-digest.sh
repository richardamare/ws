#!/usr/bin/env bash
# Render a human-readable digest of an audit log.
# Usage: .claude/hooks/audit-digest.sh [path-to.jsonl]
# Default: the most recently modified log in .claude/audit/.
root="${CLAUDE_PROJECT_DIR:-$(pwd)}"
dir="$root/.claude/audit"
file="${1:-$(ls -t "$dir"/*.jsonl 2>/dev/null | head -1)}"

if [ -z "$file" ] || [ ! -f "$file" ]; then
  echo "No audit logs in $dir"
  exit 0
fi

echo "# Audit digest: $(basename "$file")"
echo "# $(wc -l < "$file" | tr -d ' ') tool calls"
echo
echo "## Tool counts"
jq -r '.tool' "$file" | sort | uniq -c | sort -rn
echo
echo "## Timeline"
jq -r '"\(.ts)  \(.tool)  \(.summary)"' "$file"
