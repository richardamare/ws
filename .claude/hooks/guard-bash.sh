#!/usr/bin/env bash
# PreToolUse guard for Bash. Denies git ops that violate CLAUDE.md non-negotiable
# rules — including inside compound commands (`a && git commit --amend`), which the
# permissions deny-list misses because it only prefix-matches the whole command.
# Emits a PreToolUse deny decision on stdout, or nothing (allow). Always exits 0.
input=$(cat)
cmd=$(printf '%s' "$input" | jq -r '.tool_input.command // ""' 2>/dev/null)
c=$(printf '%s' "$cmd" | tr '\n' ' ')

# Strip quoted literals before matching so commit messages and grep patterns
# (e.g. a -m "...push to master..." subject) can't trigger a false git-verb match.
# Best-effort only — escaped/nested quotes can leak, so seg() does not rely on it.
scan=$(printf '%s' "$c" | sed "s/'[^']*'//g; s/\"[^\"]*\"//g")

deny() {
  jq -n --arg r "$1" '{hookSpecificOutput:{hookEventName:"PreToolUse",permissionDecision:"deny",permissionDecisionReason:$r}}'
  exit 0
}

# Match a git SUBCOMMAND: the verb must follow `git` and any global options
# (-C <dir>, -c <kv>, --git-dir/--work-tree/--namespace <v>, --exec-path=…, bare
# flags like --no-pager). This anchors to the real verb, so "push"/"amend" inside
# a quoted commit message or grep pattern never matches (no `git ` precedes it).
# $1 is the verb plus any same-segment trailing pattern (e.g. 'commit\b...--amend').
GLOBAL='([[:space:]]+(-C[[:space:]]+[^[:space:]&|;]+|-c[[:space:]]+[^[:space:]&|;]+|--git-dir[[:space:]]+[^[:space:]&|;]+|--work-tree[[:space:]]+[^[:space:]&|;]+|--namespace[[:space:]]+[^[:space:]&|;]+|--exec-path=[^[:space:]&|;]+|--no-pager|--paginate|--bare|-p))*'
seg() { printf '%s' "$scan" | grep -qE -- "(^|[^[:alnum:]_])git${GLOBAL}[[:space:]]+$1"; }

# Evaluate branch state in the repo the command actually targets: honor a leading
# `cd <dir>` or `git -C <dir>`, else the hook's cwd. Without this, a command that
# cd's into another repo is judged against the wrong branch.
repo=$(printf '%s' "$scan" | grep -oE -- '-C[[:space:]]+[^ &|;]+' | head -1 | sed -E 's/^-C[[:space:]]+//')
[ -z "$repo" ] && repo=$(printf '%s' "$scan" | grep -oE 'cd[[:space:]]+[^ &|;]+' | head -1 | sed -E 's/^cd[[:space:]]+//')
gitloc=()
[ -n "$repo" ] && [ -d "$repo" ] && gitloc=(-C "$repo")

branch=$(git "${gitloc[@]}" symbolic-ref --short HEAD 2>/dev/null)
default=$(git "${gitloc[@]}" symbolic-ref --short refs/remotes/origin/HEAD 2>/dev/null | sed 's@^origin/@@')
[ -z "$default" ] && default=master

# Rule #1 — never push to the default branch.
if seg '\bpush\b'; then
  [ "$branch" = "$default" ] && deny "Refusing git push from $default. CLAUDE.md rule #1: $default is always deployable. Branch first."
  printf '%s' "$scan" | grep -qE -- "(origin[[:space:]]+$default\b|[: ]$default\b|HEAD:$default\b)" \
    && deny "Refusing push to the $default ref. CLAUDE.md rule #1: never push to $default. Open a PR instead."
fi

# Rule #7 — no bypasses or history rewrites, even chained.
printf '%s' "$scan" | grep -qE -- '--no-verify'     && deny "Refusing --no-verify. CLAUDE.md rule #7: fix the hook failure, do not bypass it."
seg '\bcommit\b[^&|;]*--amend'                       && deny "Refusing git commit --amend. CLAUDE.md rule #7: never amend; make a new commit."
# Rule #7 narrows to the default branch: force-push to feature branches is allowed.
if seg '\bpush\b[^&|;]*(--force|[[:space:]]-f\b)'; then
  { [ "$branch" = "$default" ] \
    || printf '%s' "$scan" | grep -qE -- "(origin[[:space:]]+$default\b|[: ]$default\b|HEAD:$default\b)"; } \
    && deny "Refusing force-push to $default. CLAUDE.md rule #7: never force-push the default branch."
fi
seg '\breset\b[^&|;]*--hard'                          && deny "Refusing git reset --hard. Destroys uncommitted work; use git stash or git restore."

exit 0
