---
name: solution-design
description: |
  Produce a structured solution design document that is *the* convergence contract
  for agentic execution. The doc must be unambiguous enough that you, a future
  session, and the downstream skills (`phase-planner`, `backlog-planner`,
  `work-on-issue`) all converge on the same architecture, wire formats, schema,
  and conventions — without re-deciding settled questions.

  Triggers: "navrhni řešení", "solution design", "architektura pro", "jak bychom postavili",
  client request descriptions, pre-sales scenarios, "udělej PoC pro", any non-trivial
  multi-week build, any work that will be split across multiple agents or sessions.

  Does NOT trigger for: code questions, refactors, bug fixes, small features in existing
  code, single-PR work.
---

# Solution Design Skill

You produce a structured solution design document saved as `docs/SOLUTION.md` (or
the path the orchestrator specifies). The document is the **single point of
failure for agentic development**: if two agents read it and disagree on any
decision, downstream work diverges and drift compounds. Your job is to eliminate
ambiguity, not to write pretty prose.

## The principal-architect lens

When you write SOLUTION.md, the test you apply to every section is:

> If `work-on-issue` (or you, three weeks from now) builds from this with zero
> chance to ask a question — does it land on the same implementation, the same
> wire format, the same vocabulary you intended?

If the answer is "probably," rewrite the section until it's "yes." Implicit
defaults, "we'll figure it out later," and "use the obvious choice" are all
divergence vectors. They cost more in re-work than they save in design time.

The document is also a *contract against re-litigation*. Settled debates must be
recorded, with rejection rationale, so a future agent doesn't naively propose
what was already tried and rejected.

## Stack manifest

This is the canonical decision space. Do not propose alternatives unless the
user overrides explicitly.

### Backend variants

- **Python FastAPI (modern async)** — agents, RAG, modern async work
- **.NET 10 Azure Functions, isolated worker** — document processing, durable workflows, integrations
- **ASP.NET Core 10 Minimal APIs + EF Core** — full web apps with rich domain model (the starter default)

### Frontend variants

- **React Router 7 SPA** + Vite + **shadcn/ui owned in-repo** + Tailwind v4 + TanStack Query + Zustand + MSAL (the starter default)
- **Static HTML + minimal JS** — landing, marketing, embeddable widgets
- **Power Platform Canvas** — internal approval workflows, form-heavy admin tools

### AI / LLM variants

- **Claude Agent SDK** (Anthropic API direct) — agent loops, tool use, MCP servers
- **AWS Bedrock** (work tenant) — internal work work where Anthropic API is not procured
- **Azure OpenAI + AI Foundry** — client work where the customer is Azure-aligned
- **LangGraph** — multi-step deterministic graphs (rarely needed; default to plain Claude tool-use first)

### Data variants

- **PostgreSQL 16 Flexible Server** (+ `pgvector` when embeddings needed) — the starter default
- **Azure SQL Serverless** — when client mandates SQL Server
- **Cosmos DB DiskANN** — vector-heavy RAG with very large corpora
- **Dataverse** — Power Platform projects

### Deployment variants

- **Azure Container Apps + Static Web Apps (linked backend)** — the starter default
- **Azure Functions** — event-driven, sporadic
- **App Service** — when Container Apps is not approved at the tenant

### Always-on defaults

- **Auth:** Entra ID via MSAL (web) + in-API OAuth 2.1 proxy (MCP)
- **Observability:** OpenTelemetry → Application Insights, W3C trace context end-to-end
- **IaC:** Terraform `azurerm ~> 4.0` + `azapi ~> 2.0`, OIDC federation to GitHub Actions
- **Secrets:** Azure Key Vault, KV-referenced container secrets via user-assigned Managed Identity
- **Anthropic API model default:** latest Sonnet (`claude-sonnet-4-6` at time of writing). Never Opus or Haiku unless the user asks. Always implement prompt caching.

## Inheritance protocol

**Before writing anything**, check for parent context. If any of these exist in
the current working directory (or a starter the user references), READ them and
inherit decisions rather than re-deciding:

- `CLAUDE.md` — operational rules, conventions, "don't touch X" lists
- `docs/architecture/adr/` — Architecture Decision Records (each is load-bearing)
- `docs/patterns/` — language/framework conventions
- `docs/architecture/stack.md`, `docs/product/conventions.md`, `docs/architecture/diagrams/`
- A referenced starter repo (e.g. `work/ma.starter.web-app`)

Inherited decisions appear in your SOLUTION.md as:

> **Inherited from `work/ma.starter.web-app`:** .NET 10 + ASP.NET Core Minimal
> APIs + EF Core; React Router 7 + shadcn/ui + Tailwind v4 + TanStack Query +
> Zustand. See starter's `CLAUDE.md` and `docs/patterns/` for full conventions.

Then your stack-decisions table only enumerates **additions** and **deviations**
from the inherited baseline. This prevents the most common drift: re-deciding
things the starter already nailed, with slightly different choices.

**If no parent context is found in cwd**, ask the user once (via
`AskUserQuestion`): "Is there a starter / parent repo I should inherit
conventions from?" Options: a specific repo path, a referenced public starter,
or "greenfield — no parent." Then proceed. Never silently assume greenfield —
that's how teams end up with two slightly different React conventions in the
same org.

## Mode

Solution docs come in three modes. The mode determines depth, sections, and
elicitation rigor.

| Mode         | When                                                           | Pages | Sections                                              |
| ------------ | -------------------------------------------------------------- | ----- | ----------------------------------------------------- |
| **Vibe**     | Solo exploration, PoC for personal review, throwaway           | 1–2   | §§1, 3, 4, 5, 15 only                                 |
| **Standard** | Small team, internal feature, sales demo, single-engineer build | 3–6   | All sections, lightweight content                     |
| **Contract** | Client engagement, multi-week build, work driven through `work-on-issue` | 6–12  | All sections, *strict* convergence checklist applied  |

Detect mode from context:

- "vibe", "experiment", "rychlý", "honem" → Vibe
- "klient", "nabídka", "pre-sales", "engagement", "multi-week", agentic execution implied → **Contract** (default for unknown)
- Otherwise → Standard

**Contract is the default if you cannot tell.** The cost of over-specifying a
small project is one extra page. The cost of under-specifying an agentic build
is broken downstream work.

## Process

### Step 0 — Idempotency check

If `docs/SOLUTION.md` already exists in cwd:

1. Read it.
2. Read the change log at the bottom.
3. Do **not** overwrite. Author a new version (bump the version field), append
   a change log entry naming what changed and why, and edit in place.
4. If the user's request describes a fundamentally new system (not an evolution
   of the existing one), ask the user whether to author alongside (`docs/SOLUTION-<slug>.md`)
   or replace.

Silent overwrites destroy review history and break the contract with prior agents.

### Step 1 — Inheritance read (silent)

Run the inheritance protocol above. **Read inherited docs in full, not skim.**
A CLAUDE.md / ADR / pattern doc that you skimmed will silently re-decide the
choice it actually settled. Specifically:

- `CLAUDE.md` — every section.
- `docs/architecture/adr/*.md` — every accepted ADR. Note each ADR number for later reference in §13.
- `docs/patterns/*.md` — every pattern relevant to the stack you're designing.
- Any referenced starter's `README.md` and top-level docs.

Note what you inherit; note what you'll need to add or deviate from.

### Step 2 — Classification (silent)

Classify the request into a stack vector. Heuristics:

- "document extraction", "PDF", "kontrakt anonymizace" → .NET Functions + Document Intelligence
- "agent", "multi-step", "HITL", "pipeline" → Python + Claude Agent SDK
- "chatbot", "RAG" → Python FastAPI + React + Cosmos DiskANN
- "webhook", "integration", "trigger" → Azure Functions
- "form", "approval workflow" → consider Power Platform; otherwise React shell
- "interní pro work" + AI → AWS Bedrock
- "klient s Azure" + AI → Azure OpenAI
- "klient s M365 / Power Platform" → consider Dataverse + Power Platform
- "MCP server", "Claude Cowork" → ASP.NET Core (.NET starter) or FastAPI, with MCP SDK mounted on the existing API process — **do not split into a second service** without ADR-level justification

### Step 3 — Disambiguation (one batched question, ≤4 items)

Use `AskUserQuestion` **once**, with up to 4 questions. Ask only about
divergence-vectors that:

1. Cannot be inferred from heuristics in Step 2
2. Are not already settled by inherited CLAUDE.md / ADRs
3. Would cause measurably different downstream work depending on the answer

Sharpen each question with two or three concrete option choices.

**Always-valid questions** (when not inherited):

- Audience: internal work / client-internal / customer-facing? (drives auth choice, AI tenant)
- LLM tenant: AWS Bedrock (work) / Azure OpenAI (client) / Anthropic API direct?
- Web shape: full SPA (the starter) / minimal HTML / no web (API + MCP only)?
- Schema confidence: greenfield (you decide) / brownfield (existing DB to integrate)?

**Never ask:**

- "What framework should we use?" — pick from the manifest
- "Should we use Docker?" — yes, always
- "What language?" — pick from the manifest
- Anything inherited from CLAUDE.md / ADRs

If you cannot decide and the question is not in the always-valid set, **pick the
manifest default and mark the choice as ASSUMPTION** in §5 of the doc (Stack
decisions). Do not ask just to feel safer.

### Step 4 — Drafting

Write the full document using the template below. While drafting, every
decision goes into exactly one of three states:

- **LOCKED** — canonical. Implementation must follow.
- **DEFERRED** — explicitly punted, with named owner (ADR-XXX-planned / person / condition).
- **UNRESOLVED** — blocks execution. Listed in §14 Open questions with named owner.

There is no fourth state. "We'll figure it out" is not a state.

### Step 5 — Self-critique (convergence checklist)

Before declaring the doc done, run the **Convergence checklist** (below) against
the draft. Every line gets either a tick or a follow-up edit. Do not skip this —
it is the only thing standing between a "looks good" doc and one that survives
contact with parallel agents.

For Contract mode, the checklist is mandatory: **create a TodoWrite task per
checklist line item**, work them in order, mark each complete only when the
draft passes that line. This makes progress visible to the user and prevents
the "I checked everything mentally" rationalization.

For Standard, run it lightly. For Vibe, skip.

### Step 6 — Handoff

Output a chat message in the format at the bottom of this skill. Wait for the
user before proceeding to `phase-planner` — solution review is a human
checkpoint.

**Downstream contract.** Sections §5 (stack decisions), §6 (data model — entity
table with status markers), §7 (wire formats), §11 (testing topology), and §13
(ADR roster) are consumed by `phase-planner` and `backlog-planner`. Keep their
shapes stable — those skills grep these sections for entity names, status
markers, and ADR numbers. Renaming a section header silently breaks the chain.
The chain ends at `work-on-issue`, which executes a GitHub issue with zero
clarifying questions — so every acceptance criterion this doc sets must be
measurable enough to survive that far (`docs/product/definition-of-done.md`).

## Output template

```markdown
# Solution design: <title>

**Status:** Draft | Review | Approved
**Version:** v0.1
**Author:** ws authors (work)
**Date:** <today>
**Client / context:** <internal | client name>
**Mode:** Vibe | Standard | Contract
**Template baseline:** <starter repo or "greenfield">

> Every decision in this document is **LOCKED**, **DEFERRED** (with named ADR /
> owner / condition), or **UNRESOLVED** (listed in §14 with owner). No
> implicit defaults.
>
> **Stability:** §§2, 5, 6, 7, 8 are *contract sections* — changing them
> requires a version bump and a changelog entry. §§14, 15 evolve continuously.

### Table of contents

1. Context
2. Glossary
3. Requirements
4. Architecture
5. Stack decisions
6. Data model
7. Wire formats & API conventions
8. Security & compliance
9. Observability
10. Deployment
11. Concurrency, lifecycle, and operational conventions
12. Documentation plan
13. ADR roster
14. Risks & open questions
15. Estimation

---

## 1. Context

Why does this exist? What business problem does it solve? Who is the user?
2–3 paragraphs max. End with a one-sentence value proposition that an agent
can quote.

## 2. Glossary

Define every domain term used in the rest of the document. One line per term.
Agents reading the doc must use this vocabulary verbatim — no synonyms.

| Term | Definition |
| ---- | ---------- |
| Scenario | A discrete AI deployment opportunity for one customer, with scoring, lifecycle state, and provenance. |
| Draft | A proposed mutation from Claude/MCP awaiting human confirmation. Never canonical. |
| (etc.) | |

**Skip only in Vibe mode.** Standard and Contract require this.

## 3. Requirements

### Functional (numbered, one line each)
| #  | Module |
| -- | ------ |
| M1 | |

### Non-functional
- **Performance:** quantified targets (p95 latency, throughput, page-load budget)
- **Scale envelope:** user count, data volume, single/multi-region
- **Security:** auth flows, PII handling, secret storage
- **Compliance:** GDPR, EU AI Act position, retention windows
- **Observability:** trace propagation boundaries, alert SLOs
- **Data residency:** region + reason

### Out of scope (explicitly enumerated)
- (Things a reasonable reader would assume in scope but are NOT — protects against scope creep and silent assumption.)

## 4. Architecture

### Component diagram (Mermaid C4-container)
Show every service, transport, identity store, observability sink, and data
store. If a component appears in §5 (data) or §6 (security), it must appear
here.

```mermaid
flowchart LR
  ...
```

### Data flows (one per critical use-case)
1. Sequence diagram per primary flow (auth, primary write, primary read, AI flow).
2. Show actor → component → component → data store, with transport labels.

### Integration points
- External system, transport, identity, ownership, fallback when unavailable.

## 5. Stack decisions

> Decisions inherited from the parent starter are listed first as "Inherited";
> only additions and deviations carry a "Why" column.

| Concern | Choice | Status | Why |
| ------- | ------ | ------ | --- |
| Backend | .NET 10 + ASP.NET Core Minimal APIs + EF Core | Inherited | Starter default; see starter `CLAUDE.md`. |
| Frontend | React Router 7 + shadcn/ui (owned in-repo) + Tailwind v4 + TanStack Query + Zustand | Inherited | Starter default; UI library is shadcn/ui — **not Base UI, not Radix raw**. |
| LLM | Claude Cowork (external) | LOCKED | All AI activity in Cowork; app is deterministic system of record. Avoids EU AI Act high-risk classification. |
| Data | PostgreSQL 16 + pgvector | LOCKED (deviation: pgvector via `azure.extensions = VECTOR`) | Starter default + embedding table for template match. |
| Background jobs | Hangfire on Postgres, own schema | LOCKED (addition) | Required for draft TTL expiry, audit purge. **Not** in starter. No separate worker container. |
| MCP transport | `ModelContextProtocol.AspNetCore` at `/mcp` | LOCKED | Same process as HTTP API. Fallback: hand-rolled SSE (~1 MD). |
| Excel I/O | ClosedXML | LOCKED (addition) | .NET equivalent of openpyxl. |

### Considered and rejected

This is the *agent-defensive* section. List things an agent might naively
propose that have already been tried, considered, or ruled out — with the
reason. A future agent must not re-litigate these without an ADR.

*(Examples below are drawn from a prior MCP-on-Postgres build. Substitute for
your domain. The point is the shape: every entry is a temptation + reason it
was rejected, so a future agent can't re-propose it innocently.)*

- **Two services (HTTP API + MCP worker).** Same DB, same auth, same logic — one process is cheaper.
- **Power Platform front-end.** Custom UI required for diff-based draft confirmation; Canvas can't host this well.
- **In-app AI (risk scoring, recommendations).** Cowork has M365 context; duplicates badly. Also pulls the app under EU AI Act high-risk scope.
- **Repository pattern.** Starter convention is direct `DbContext` in services. (Was tried and removed in starter — see the relevant ADR if present.)
- **Strongly-typed ID structs** (e.g. `record struct ScenarioId`). Tried and removed; plain `Guid` on wire. Do not reintroduce without an ADR.

## 6. Data model

Schema is the most divergent surface across parallel agents. Every entity must
be enumerated, with status, key fields, and relationships. Fields not listed
here are *not* part of the contract — agents must propose them via the change
process, not invent them.

For each entity, mark:

- **Status:** MVP / deferred / reserved-table-only
- **Owner:** which feature module
- **Concurrency:** `row_version` (xmin) / append-only / single-writer

*(Example below is illustrative — substitute your real entities. The point is
the shape: status marker on every entity, owner module named, concurrency
strategy declared.)*

```
Customer (MVP) — owned by M1
  id, name, sector, sharepoint_url?, country (ISO-3166-1, default 'CZ'),
  default_currency (ISO-4217, default 'CZK'),
  status [active | archived],
  created_at, created_by, updated_at, updated_by, row_version
  ├── UserCustomerAccess (MVP, M:N, role [Admin | Manager | Consultant | Viewer])
  ├── ScoringConfig (MVP, 1:1)
  └── Scenario (MVP, 1:N)
        ├── lineage:    template_id (nullable FK → Template)
        ├── tags:       Tag (M:N global)
        ├── core:       name, author, description, ...
        ├── computed_stored: priority_score, quadrant, tco_1y, tco_3y, save_1y, save_3y
        ├── workflow:   state [Idea | Analysis | Approved | InProgress | Operating | Archived | Rejected]
        ├── concurrency: row_version (xmin)
        └── audit:      created_at, created_by, updated_at, updated_by

DraftMutation (MVP) — owned by M6
  field-level proposals. status [pending | applied | rejected | superseded | expired].
  base_row_version captured at proposal time; latest-wins on apply.

AuditLog (MVP) — append-only, retention ≥ 90 days
StateTransition (MVP) — append-only, retention ≥ 365 days
SystemConfig (MVP) — runtime knobs (draft TTL hours, etc.); DB-backed, not appsettings.

ScenarioEmbedding (DEFERRED to Phase 2) — vector(3072), pgvector
TemplateEmbedding (DEFERRED to Phase 2) — vector(3072), pgvector
```

### Deletion semantics (per entity)

| Entity | Soft / Hard | Default rule | Admin escape hatch |
| ------ | ----------- | ------------ | ------------------ |
| Scenario | Archive then hard-delete | Cannot hard-delete with pending drafts | `?hard=true` admin-only, audited |
| Customer | Archive then hard-delete | Cascade scenarios on hard-delete | `?hard=true` admin-only |
| AuditLog | Never deleted by hand | Hangfire `audit-purge` past retention | None |

### Design choices (load-bearing)

- **Computed fields are persisted, not derived on read.** Dashboards stay fast; historical values survive scoring-config changes.
- **`DraftMutation` is field-level**, not whole-record. N changes from Claude = N rows. Diff UI becomes trivial.
- **No soft delete column.** `Archived` state covers the use case.
- **`row_version` mapped to Postgres `xmin`** for optimistic concurrency on PUT/PATCH and draft-apply.

## 7. Wire formats & API conventions

This section exists because two agents will pick differently if you don't lock
this. Every line is a contract.

| Concern | Lock |
| ------- | ---- |
| JSON field naming | `camelCase` on the wire (API serialization policy + matching Zod schemas on web). |
| ID format on the wire | Plain `Guid` (lowercase string), never branded/strongly-typed. (Tried and reverted — see Considered and rejected.) |
| Enum format on the wire | `snake_case` string, never numeric. (`PropertyNamingPolicy` + custom enum converter on API; Zod string union on web.) Unknown enum value on input → `400 validation.failed`. |
| Date format | ISO-8601 with UTC `Z` suffix, never with offset. Date-only fields use `YYYY-MM-DD` (no time, no zone). |
| Null vs omit | Optional fields are omitted when not set, never sent as explicit `null`. Receivers treat omit and `null` identically. |
| Booleans | `true` / `false` only. Never `0` / `1`, never `"true"` / `"false"`. |
| Money | Integer minor units (e.g. CZK halíře, EUR cents) + ISO-4217 currency code in a sibling field. Never floats. |
| Pagination — list | `?skip=<int>&take=<int>`, response includes `{ items, total }`. |
| Pagination — search | Cursor: `?cursor=<opaque>&take=<int>`, response includes `{ items, nextCursor }`. |
| Error response | RFC 7807 ProblemDetails with stable `code` field. MCP surfaces same `code` via JSON-RPC `error.data.code` (numeric `-32001` for business errors). |
| Concurrency | `If-Match: <row_version>` header on PUT/PATCH; mismatch → `409 { code: "<entity>.concurrency_conflict", currentRowVersion: "<x>" }`. |
| Auth header | `Authorization: Bearer <jwt>`. Token from MSAL (web) or in-API OAuth proxy (MCP). |
| CORS | Allowlist of SWA hostnames only, per env. |
| Trace propagation | W3C `traceparent` header across web → API → MCP. |
| Cancellation | Every API endpoint takes `CancellationToken`; long-running ops respect it. |

### Error code catalog (excerpt; full catalog lives in source)

`auth.unauthenticated`, `auth.forbidden`, `customer.not_found`,
`scenario.not_found`, `scenario.archived`, `scenario.concurrency_conflict`,
`draft.not_found`, `draft.expired`, `draft.superseded`,
`draft.duplicate_import_row`, `validation.failed`, `template.not_found`,
`tag.namespace_locked`, `excel.parse_failed`, `rate_limit.exceeded`,
`payload.too_large`, `internal.unexpected`.

## 8. Security & compliance

- **Auth (web):** MSAL flow, token cache lives in MSAL (no Zustand mirror), bearer injected per request.
- **Auth (MCP):** Entra OAuth 2.1 brokered by in-API proxy at `/oauth/*` + `/.well-known/oauth-*`. JWT validated against v1 *and* v2 issuers (v1 accepted to support desktop OAuth's `resource` parameter). Required scope: `<scope-name>`.
- **Identity in audit log:** `actor_user_id` resolved from JWT `oid` claim; `source = mcp` distinguishes MCP-originated mutations.
- **Authorization:** per-customer roles (`Admin`, `Manager`, `Consultant`, `Viewer`) plus global `User.is_admin`. MCP inherits authenticated user's roles — no separate key scope.
- **Secrets:** Key Vault (RBAC), KV-referenced container secrets via user-assigned Managed Identity. No plain env-var secrets.
- **PII inventory:** enumerate every field that holds PII, why, retention, deletion path.
- **EU AI Act position:** classify (high-risk / limited-risk / out-of-scope) with reasoning. Document in the work AI register.
- **Platform-level auth (Easy Auth / SWA auth):** disabled. Auth is in-app.

## 9. Observability

- **OTel SDK:** instrumentations (ASP.NET Core, HttpClient, EF Core / Npgsql).
- **Logging:** Serilog → Application Insights connection string from KV. Request enrichment: `traceId`, `actor_user_id`, `customer_id`, `source`.
- **Noise filters:** explicitly enumerate spans that must be dropped (e.g. Hangfire poll loop parentless Npgsql spans dropped at `OnStart` before Azure Monitor batch exporter). New high-frequency parentless sources MUST extend this list.
- **Metrics:** name each custom metric (request duration histogram, draft lifecycle counter, etc.).
- **Alerts (prod):** thresholds, evaluation windows, routing target.

## 10. Deployment

### Topology per environment

Enumerate every Azure resource per env (`dev`, `prod`), with SKU, region, and
naming pattern.

### IaC layout

Folder structure under `deploy/terraform/`, backend (state) configuration, how
environments are parameterized. Note any places where `azapi` is required because
`azurerm` lags.

### CI/CD

Per-pipeline list: trigger, steps, gates (typecheck, format, integration tests
with real Postgres via Testcontainers, OpenAPI drift check), deployment target.
Explicit list of CI-enforced rules:

- `bun run gen:api:check` — no drift between API DTOs and generated web client
- `dotnet format --verify-no-changes`
- `terraform validate`

### Cost envelope (EUR / month)

| Resource | Dev | Prod |
| -------- | --- | ---- |
|          |     |      |
| **Total** |    |      |

### Rollout (brownfield only — skip section if greenfield)

If this system replaces or augments an existing one:

- **Current state:** what exists today, where data lives, who uses it.
- **Coexistence window:** dual-write, read-from-old / write-to-new, or hard cutover?
- **Migration path:** one-shot import / continuous sync / per-entity backfill?
- **Cutover criteria:** measurable conditions that gate the switch.
- **Rollback plan:** how to revert if the new system fails post-cutover.

## 11. Concurrency, lifecycle, and operational conventions

- **Optimistic concurrency:** scope, header, response on mismatch.
- **Draft lifecycle:** state machine, TTL, supersede rules, archived-collapse rules.
- **Background jobs:** name, schedule (UTC), purpose, idempotency claim.
- **Rate limits:** per-user, per-tool-class, per-payload size, response shape on exceed.
- **Migrations:** how to add one (exact `dotnet ef migrations add` invocation), naming, when to data-migrate vs schema-migrate.
- **API contract change protocol:** "DTO/endpoint/enum change → `bun run gen:api` in the same commit; CI enforces drift-free."
- **Local dev:** one command to bring up dependencies (`docker compose -f apps/api/docker-compose.yml up -d`), env vars required, where seed data lives.

### Testing topology

Tests are a divergence vector if not specified. Lock the topology:

| Test project | Scope | Real / mocked DB | When it runs |
| ------------ | ----- | ---------------- | ------------ |
| `tests/AiRadar.Api.Tests` | HTTP + MCP endpoint integration | Real Postgres via Testcontainers — **never mocked** | CI on every PR; locally with Docker |
| `tests/AiRadar.Scoring.Tests` | Pure scoring engine | n/a (no DB) | CI + locally; fast |
| `apps/web/e2e` (Playwright) | Browser-level happy paths | Real backend (`dev` or local) | CI on merge to `master` |

- **Seed data:** authoritative seeder under `apps/api/src/<project>/Persistence/DataSeeder.cs`; test-only fixtures under `tests/<project>/Fixtures/`.
- **Mocking rule:** mock external HTTP boundaries (Entra, AI Foundry) only. Do not mock the database, do not mock the MCP transport.

## 12. Documentation plan

The codebase will need supporting docs beyond this one. Enumerate them so an
agent knows what to write when shipping a feature.

| Doc | Purpose | Owner |
| --- | ------- | ----- |
| `CLAUDE.md` | Operational rules for Claude in this repo | Architect |
| `docs/architecture/adr/ADR-NNN-*.md` | One per load-bearing decision (see ADR roster §13) | Per decision |
| `docs/patterns/csharp.md` | C# conventions (nullability, records, file layout) | Lead engineer |
| `docs/patterns/react.md` | React conventions (server vs UI state, forms, auth gate) | Lead engineer |
| `docs/patterns/terraform.md` | IaC conventions | Lead engineer |
| `docs/runbooks/<name>.md` | One per operational procedure (deploy, rollback, secret rotation) | Per phase |
| `.claude/skills/<skill>/SKILL.md` | Project-specific Claude skills (e.g. domain-aware propose helpers) | Per skill |

## 13. ADR roster

ADRs that will be authored during build. Numbered preemptively so commits can
reference them. Each is one decision; each is load-bearing.

| ADR | Title | Trigger |
| --- | ----- | ------- |
| ADR-001 | Platform: Container Apps + Static Web Apps linked backend | Phase 1 |
| ADR-002 | Database: PostgreSQL 16 + pgvector | Phase 1 |
| ADR-003 | Auth: Entra ID v1+v2 issuers, in-API OAuth proxy for MCP | Phase 2 |
| ADR-004 | MCP transport: `ModelContextProtocol.AspNetCore` at `/mcp` | Phase 2 |
| ADR-005 | Drafts: field-level, propose-only MCP surface | Phase 2 |
| ADR-006 | Scoring: deterministic, persisted, server-side only | Phase 1 |
| ADR-007 | IDs: plain `Guid` on wire — strongly-typed structs rejected | Phase 1 |
| ADR-008 | Pagination: skip/take for lists, cursor for search | Phase 1 |
| ADR-009 | Errors: `AppException` + `ErrorCodes` → RFC 7807 ProblemDetails | Phase 1 |
| ADR-010 | Enums: snake_case strings on the wire | Phase 1 |

## 14. Risks & open questions

### Risk register

| Risk | Likelihood | Impact | Mitigation | Owner |
| ---- | ---------- | ------ | ---------- | ----- |

### Open questions (UNRESOLVED — block execution)

| # | Question | Blocks | Owner | Due |
| - | -------- | ------ | ----- | --- |

### Deferred decisions (DEFERRED — execution proceeds; revisit at trigger)

| # | Decision deferred | Trigger to revisit | Owner |
| - | ----------------- | ------------------ | ----- |

## 15. Estimation

| Phase | Scope | Effort (MD) |
| ----- | ----- | ----------- |

**Total:** X–Y MD over <calendar window>.
**Confidence:** High / Medium / Low — with reasoning (what removes uncertainty, what residual unknowns remain).

---

**Change log**

| Version | Date | Author | Change |
| ------- | ---- | ------ | ------ |
| v0.1    |      |        | Initial draft |
```

## Convergence checklist

Run this before declaring the doc done. Each line gets a tick. If any line
fails, edit and re-run.

**Vocabulary**
- [ ] §2 Glossary defines every non-obvious domain term used in §§3–14.
- [ ] No synonym for a glossary term appears elsewhere (search the doc).

**Stack**
- [ ] Every row in §5 Stack decisions has a Status (Inherited / LOCKED / DEFERRED / UNRESOLVED).
- [ ] The UI component library is named explicitly (e.g. "shadcn/ui owned in-repo"). Not "a React UI library."
- [ ] Anthropic API model (if any) is pinned to a specific Sonnet version, with prompt caching mandated.
- [ ] §5 "Considered and rejected" lists at least three temptations an agent might naively propose.

**Schema**
- [ ] Every entity in §6 has Status (MVP / deferred / reserved-table-only).
- [ ] Aspirational entities (e.g. embeddings before they're built) are marked DEFERRED, not listed as MVP.
- [ ] Concurrency strategy is named per entity.
- [ ] Deletion semantics enumerated per entity, including admin escape hatches.

**Wire formats**
- [ ] §7 locks JSON field naming, ID format, enum casing, date format, null-vs-omit, booleans, money, pagination shape, error shape, concurrency header.
- [ ] Error code catalog excerpt present and stable (codes are part of the contract).

**Testing**
- [ ] §11 names every test project, its scope, what's real vs mocked, when it runs.
- [ ] Mocking rule is explicit (e.g. "do not mock the DB").

**Security**
- [ ] Auth flow for every actor (web user, MCP client, background job) is specified.
- [ ] JWT claim used for `actor_user_id` is named.
- [ ] PII inventory is concrete (which fields, why, retention).
- [ ] EU AI Act position is stated with reasoning.

**Observability**
- [ ] Noise filters explicitly enumerated.
- [ ] Custom metrics named.

**Deployment**
- [ ] Every Azure resource enumerated per env, with SKU and region.
- [ ] CI-enforced rules listed (drift check, format, tests).
- [ ] Cost envelope is per-env, with line items.

**Process / doc plan**
- [ ] §12 lists all supporting docs the codebase will need.
- [ ] §13 names planned ADRs preemptively.
- [ ] §11 includes the API contract change protocol verbatim.

**Decision state**
- [ ] Every decision in the doc is LOCKED, DEFERRED (with owner+trigger), or UNRESOLVED (in §14, with owner+due).
- [ ] §14 Open questions has named owners and due dates, not "TBD".

**Versioning**
- [ ] Top of doc has Version + change log at bottom.

## Anti-patterns

- **Pretty prose over precision.** "We'll use a modern UI library" is a divergence vector. Name the library.
- **Re-deciding inherited choices.** If the starter's CLAUDE.md says shadcn/ui, do not write "React component library: TBD" in your SOLUTION.
- **Aspirational schema.** Entities you might build are not entities you will build. Mark them DEFERRED.
- **Silent defaults.** "Standard pagination" / "the usual error model" — agents will pick differently. Lock it in §7.
- **"We'll figure it out."** Not a decision state. Use DEFERRED with a trigger, or UNRESOLVED with an owner.
- **Open questions without owners.** An open question without a name attached is a permanent open question.
- **Skipping the glossary in non-Vibe mode.** Agents that don't share vocabulary don't converge.
- **Asking more than 4 disambiguation questions.** If you need more, you're missing inherited context — go re-read CLAUDE.md.
- **Writing code in SOLUTION.md.** It's a WHAT contract, not a HOW reference.
- **Proceeding to `phase-planner` automatically.** Solution review is a human checkpoint. Always.
- **Treating SOLUTION.md as immutable after v1.** It evolves. Bump the version, log the change, link to the ADR that drove it.

## Handoff message

After writing the doc, post:

```
SOLUTION.md vytvořen (v0.1, mode: <Vibe|Standard|Contract>).

Convergence checklist: <PASS — N/N> | <FAIL — list failed line items>

Klíčová LOCKED rozhodnutí:
1. <stack/architecture decision 1, one line>
2. <decision 2>
3. <decision 3>

DEFERRED (execution běží, revisit při triggeru):
- <decision 1> → ADR-<n>, trigger: <when>
- <decision 2> → <person>, trigger: <when>

UNRESOLVED (blokuje execution) — <N> položek vyžaduje rozhodnutí PŘED `phase-planner`:
1. <question 1> — owner: <name>
2. <question 2> — owner: <name>

Bez rozhodnutí každé UNRESOLVED položky se downstream agenti rozejdou. Jakou volíš:
  (a) Resolve teď — řekni mi rozhodnutí, doplním do §5/§6/§7 jako LOCKED.
  (b) Defer — pojmenuj ownera + trigger, přesunu do "Deferred decisions".
  (c) Accept block — UNRESOLVED zůstává, `phase-planner` neběží dokud se to nevyřeší.

Convergence checklist failures (pokud nějaké):
- <failed item> — <one-line reason>

Next: rozhodni o UNRESOLVED → fix checklist failures → `phase-planner`.
```

If there are zero UNRESOLVED items AND zero checklist failures, the message
ends at "DEFERRED" and the next-step is simply `phase-planner`. Otherwise the
skill **must not** silently proceed.

## Length scaling

- **Vibe:** §§1, 3, 4, 5, 15 only. Skip glossary (§2), wire formats (§7), doc plan (§12), ADR roster (§13), convergence checklist.
- **Standard:** all sections, lightweight content (3–6 pages). Convergence checklist run lightly.
- **Contract:** all sections, full content (6–12 pages). Convergence checklist mandatory. Every checklist failure is a blocking edit.

Default to **Contract** when in doubt. Over-specifying a small project costs a
page; under-specifying an agentic build costs a week.
