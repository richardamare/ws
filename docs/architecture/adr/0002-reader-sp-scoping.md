# ADR-0002: Scope sessions with a per-project Reader service principal

- **Status:** Accepted
- **Date:** 2026-06-27
- **Deciders:** ws authors

## Context

the maintainer holds near-Owner (Contributor + User Access Administrator + RBAC Administrator), permanent and
inherited at **subscription** scope, on a **shared lab subscription**. Any session logged in as him can
delete any resource group in the lab. The permissions are slow/bureaucratic to regain, so removing them
or converting them to PIM-eligible is unacceptable. A session cannot reduce its own token's permissions —
Azure RBAC is identity-based. User-assigned managed identities do not authenticate from local Docker or
the host Mac (no IMDS endpoint).

## Decision

We will have everyday `ws` sessions authenticate as a **separate service principal with Reader on a
single resource group** — one SP per project (`sp-ws-<proj>-reader`) — isolated from the personal
admin login via a per-project **`AZURE_CONFIG_DIR`**. The personal admin identity is never modified and
is used only via the deliberate `ws elevate` path (write/Terraform, human-approved).

## Consequences

- Easier: default sessions are read-only on one RG; `az group delete` elsewhere → AuthorizationFailed.
  the maintainer keeps all his roles, untouched.
- Harder / accepted: write requires a deliberate personal login (friction by design); one machine
  identity per project must be created and cleaned up; certs expire (~1yr) → `ws rotate`.
- Follow-on: defense-in-depth via RG `CanNotDelete` locks and the `.claude/settings.json` deny list.
- Ruled out: PIM-eligible (fear of losing standing access) and UAMI (no local IMDS).
