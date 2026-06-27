# Security — Azure trust model

**The most important doc in this repo. Read before writing any code that touches `az`.**

## The problem

The maintainer has **near-Owner** on a **shared lab subscription**: Contributor + User Access Administrator +
RBAC Administrator, all **permanent-active**, inherited at **subscription** scope (0 PIM-eligible,
0 deny assignments). So today his identity can delete **any** resource group in the lab — and any
session logged in as him inherits that. A stray/hallucinated `az group delete` is a real risk.

These permissions are slow/bureaucratic to regain, so **they must not be removed or converted to
PIM-eligible.** The account stays exactly as is.

## The principle

You cannot shrink your own token — Azure ties permissions to identity. So the everyday session
authenticates as a **different, Reader-only principal**, while the personal admin identity sits
untouched in a separate credential store and is used only deliberately.

**Scope the session, not the account.**

## The mechanism — two identities, isolated by `AZURE_CONFIG_DIR`

- **Personal admin** → normal `~/.azure` — full roles, untouched. Used deliberately for write work.
- **Everyday `ws` sessions** → `AZURE_CONFIG_DIR=~/.azure-<project>`, logged in as the Reader SP.
  A session pointed here literally cannot see or use the admin token.

## The Reader principal

- A **service principal** (machine identity — NOT a new user account; no login/mailbox/license). Same
  kind of identity GitHub Actions already uses.
- **One per project**, each with **Reader** on **one resource group only**.
- Naming: `sp-ws-<project>-reader` (the `sp-ws` prefix makes all of them filterable in Entra).

```bash
az ad sp create-for-rbac \
  --name sp-ws-proj1-reader \
  --role Reader \
  --scopes /subscriptions/<sub>/resourceGroups/<rg1> \
  --create-cert
```

Result: `az group delete -n any-other-rg` from a `ws` session → **AuthorizationFailed**. Intended blast
radius: read-only, one RG.

## The write path — deliberate, separate

The SP is **Reader-only on purpose.** When write is needed (e.g. Terraform bootstrap, which the CI
managed identity is intentionally too weak to do), `ws` does **not** use the SP. `ws elevate` opens a
clearly-marked elevated tab that runs `az login` as the maintainer (`unset AZURE_CONFIG_DIR` → personal
`~/.azure`); he reviews `terraform plan`, approves `apply`, then closes the tab. Full power, only when
asked, only while watched. `terraform apply` is never automated.

## Rejected alternatives

- **User-assigned managed identity** only authenticates from Azure-hosted compute (IMDS at
  169.254.169.254). A local Docker container / the host Mac cannot use one. UAMI stays the right choice
  for CI (GitHub Actions), not for local sessions.
- **PIM-eligible** would work but the maintainer fears losing standing access and the slow re-grant. The
  separate-Reader-SP approach achieves least-privilege-by-default **without touching existing
  assignments**, so it is preferred.

## Defense-in-depth (NOT the boundary — RBAC is)

- `CanNotDelete` **lock** on important RGs (the Reader SP can't remove it; it lacks UAA).
- `.claude/settings.json` denies `az group delete`, `az role assignment create/delete`, `terraform
  apply/destroy`, `az ad sp delete`.
- Cert files `chmod 600`; SP certs expire (~1yr) → `ws rotate`.

## Entra diagnosis (role assignments / app registrations)

These are **Microsoft Graph / directory** permissions, NOT Azure RBAC — an RG-scoped identity has none
by design. Keep them OUT of the Reader SP. When needed, use the personal admin identity, deliberately.

## Hard rules (also in root CLAUDE.md)

1. Never modify the maintainer's existing Azure role assignments.
2. The per-project SP is Reader-only; never grant it write, never use it for write.
3. Write only via deliberate personal `az login` (`ws elevate`); `terraform apply` needs human approval.
