#!/usr/bin/env bash
#
# setup-repo-settings.sh – apply a standard GitHub repo configuration
# (settings + rulesets) idempotently.
#
# Re-running is safe: repo settings are PATCHed and rulesets are matched by
# name, deleted, then recreated.
#
# Requires: gh CLI authenticated with admin on the target repo, plus jq.
#
# Usage:
#   REPO=owner/name ./scripts/setup-repo-settings.sh           # apply
#   REPO=owner/name DRY_RUN=1 ./scripts/setup-repo-settings.sh # print only

set -euo pipefail

REPO="${REPO:?set REPO=owner/name}"
DRY_RUN="${DRY_RUN:-0}"

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
run()  { if [[ "$DRY_RUN" == "1" ]]; then echo "DRY: $*"; else "$@"; fi; }

for bin in gh jq; do
  command -v "$bin" >/dev/null || { echo "missing dependency: $bin" >&2; exit 1; }
done

gh api "repos/$REPO" >/dev/null 2>&1 || {
  echo "cannot access repos/$REPO – check auth and admin rights" >&2; exit 1; }

# ---------------------------------------------------------------------------
# 1. Repository settings (merge strategy, branch cleanup, features)
# ---------------------------------------------------------------------------
log "Applying repository settings to $REPO"
run gh api -X PATCH "repos/$REPO" \
  -F allow_squash_merge=true \
  -F allow_merge_commit=false \
  -F allow_rebase_merge=false \
  -F allow_auto_merge=true \
  -F delete_branch_on_merge=true \
  -F allow_update_branch=false \
  -F use_squash_pr_title_as_default=true \
  -f squash_merge_commit_title=PR_TITLE \
  -f squash_merge_commit_message=PR_BODY \
  -F has_issues=true \
  -F has_projects=true \
  -F has_wiki=false \
  -F web_commit_signoff_required=false \
  >/dev/null

# ---------------------------------------------------------------------------
# 2. Rulesets – recreate idempotently (delete any existing match by name)
# ---------------------------------------------------------------------------
delete_ruleset_by_name() {
  local name="$1"
  gh api "repos/$REPO/rulesets" --jq ".[] | select(.name == \"$name\") | .id" \
    | while read -r id; do
        [[ -n "$id" ]] || continue
        log "Deleting existing ruleset '$name' (id $id)"
        run gh api -X DELETE "repos/$REPO/rulesets/$id" >/dev/null
      done
}

create_ruleset() {
  local name="$1" body="$2"
  log "Creating ruleset '$name'"
  if [[ "$DRY_RUN" == "1" ]]; then echo "DRY: POST repos/$REPO/rulesets <<< $body"; return; fi
  echo "$body" | gh api -X POST "repos/$REPO/rulesets" --input - >/dev/null
}

# bypass_actors actor_id 5 = built-in "Admin" RepositoryRole, bypass always.
ADMIN_BYPASS='[{"actor_id":5,"actor_type":"RepositoryRole","bypass_mode":"always"}]'

# 2a. Branch ruleset on default branch: squash-only PR, linear history, no force-push, no deletion
BRANCH_NAME="master: squash-only PR, linear history, no force-push"
delete_ruleset_by_name "$BRANCH_NAME"
create_ruleset "$BRANCH_NAME" "$(jq -n --argjson bypass "$ADMIN_BYPASS" '{
  name: "master: squash-only PR, linear history, no force-push",
  target: "branch",
  enforcement: "active",
  conditions: { ref_name: { include: ["~DEFAULT_BRANCH"], exclude: [] } },
  bypass_actors: $bypass,
  rules: [
    { type: "deletion" },
    { type: "non_fast_forward" },
    { type: "required_linear_history" },
    { type: "pull_request", parameters: {
        required_approving_review_count: 0,
        dismiss_stale_reviews_on_push: false,
        required_reviewers: [],
        require_code_owner_review: false,
        require_last_push_approval: false,
        required_review_thread_resolution: false,
        allowed_merge_methods: ["squash"]
    } }
  ]
}')"

# 2b. Tag ruleset on v*: no delete, no force-update
TAG_NAME="release tags v*: no delete, no force-update"
delete_ruleset_by_name "$TAG_NAME"
create_ruleset "$TAG_NAME" "$(jq -n --argjson bypass "$ADMIN_BYPASS" '{
  name: "release tags v*: no delete, no force-update",
  target: "tag",
  enforcement: "active",
  conditions: { ref_name: { include: ["refs/tags/v*"], exclude: [] } },
  bypass_actors: $bypass,
  rules: [
    { type: "deletion" },
    { type: "non_fast_forward" }
  ]
}')"

log "Done. Current rulesets:"
run gh api "repos/$REPO/rulesets" --jq '.[] | "  - \(.name) [\(.target), \(.enforcement)]"'
