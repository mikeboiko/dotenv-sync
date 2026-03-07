<!--
Sync Impact Report
Version change: 1.0.0 -> 1.0.1
Modified principles:
- None
Added sections:
- None
Removed sections:
- None
Templates requiring updates:
- ✅ updated: .specify/templates/plan-template.md
- ✅ updated: .specify/templates/tasks-template.md
- ✅ reviewed: .specify/templates/commands/ (directory created; no command templates yet)
Runtime guidance updates:
- ✅ updated: .github/agents/speckit.constitution.agent.md
- ✅ updated: .github/agents/speckit.plan.agent.md
- ✅ updated: .github/agents/speckit.specify.agent.md
- ✅ updated: .github/agents/speckit.tasks.agent.md
Derived guidance updates:
- ✅ added: .specify/memory/constitution-checks.json
Follow-up TODOs:
- None
-->

# dotenv-sync Constitution

## Core Principles

### I. Workflow-Preserving Architecture

`dotenv-sync` MUST preserve the standard `.env` workflow. Features MUST write,
diff, or validate local files without runtime command wrappers, shell shims, or
process injection, and `.env.example` MUST remain the schema contract for
expected variables. Code changes MUST keep CLI routing, envfile parsing/writing,
sync orchestration, and provider adapters in focused modules with explicit
interfaces rather than cross-cutting logic. Rationale: adoption depends on
fitting existing developer habits, and clear boundaries keep secret-handling
code reviewable.

### II. Deterministic Env File Fidelity

Parsing and writing MUST preserve key order, comments, and unchanged lines
whenever the source format allows it. Commands MUST avoid rewriting `.env` or
`.env.example` when the semantic content is unchanged, and dry-run/diff output
MUST describe only real mutations. Reverse sync MUST add schema entries as blank
placeholders and MUST NOT synthesize or expose secret values. Rationale:
deterministic output keeps reviews clean and makes repeated sync operations safe
to trust.

### III. Test-First Reliability

Every behavior change MUST begin with failing automated tests at the narrowest
useful level, and work is not complete until unit plus CLI/integration coverage
verifies forward sync, reverse sync, validation, and provider failure paths
affected by the change. Secret-handling flows MUST include negative-path tests
for missing CLIs, locked vaults, unmapped keys, and redaction rules. Rationale:
this tool coordinates local files and external secret systems, so regressions
are expensive and often invisible without automation.

### IV. Consistent UX and Secret-Safe Output

All commands MUST use consistent status terms, exit-code semantics, and recovery
guidance across success, dry-run, validation, and failure paths. User-visible
output MUST never print raw secret values; logs, previews, diffs, diagnostics,
and errors MUST use redacted or presence-only indicators. When prerequisites are
missing, the CLI MUST state the problem, the impact, and the next operator
action. Rationale: predictable UX and strict redaction are necessary for trust
in a tool that touches secrets.

### V. Performance Budgets and Minimal Dependencies

Local parse, diff, validate, and no-op sync detection MUST complete within
200 ms p95 for `.env`/`.env.example` pairs up to 500 keys when provider
responses are mocked or cached. Provider-backed execution MUST perform no more
than one lookup per distinct key per command unless the provider API requires
pagination or re-authentication, and any new dependency MUST be justified in the
plan as a material reduction in complexity or risk. Rationale: CLI tools win on
fast feedback, low installation friction, and predictable maintenance cost.

## Product Constraints

- The implementation MUST remain a Go CLI with a minimal dependency footprint.
  The standard library or small, well-maintained packages MUST be preferred
  unless a larger dependency materially reduces complexity or risk.
- `.env.example` is the canonical schema. Empty values signal provider
  resolution, explicit values are copied verbatim, and reverse sync may only add
  blank schema placeholders to `.env.example`.
- Secret providers MUST implement a shared interface and keep provider-specific
  behavior isolated under provider adapter packages.
- The default onboarding path MUST be zero-config when `.env.example` and a
  default provider are detectable, including actionable diagnostics for missing
  CLIs, locked sessions, or unavailable vaults.
- Commands that read, compare, or display secret-bearing values MUST redact
  them unless they are writing the developer's local `.env` file.

## Delivery Workflow & Quality Gates

- Plans MUST pass a constitution check covering workflow preservation, module
  boundaries, deterministic file handling, required automated tests, UX
  consistency, secret redaction, and performance/dependency budgets before
  design work proceeds.
- Specifications MUST define measurable acceptance criteria for user-visible
  behavior, dry-run/diff trustworthiness, error recovery, and performance
  expectations for affected flows.
- Tasks MUST include test tasks before implementation plus any UX copy,
  redaction, and performance validation work needed to satisfy the constitution.
- Reviews MUST block merges when a change weakens deterministic file behavior,
  leaks secrets, omits failing-first tests, or introduces unjustified
  dependency or latency costs.
- Documentation and help output MUST be updated in the same change whenever a
  command, flag, provider prerequisite, or operator-visible message changes.

## Governance

This constitution overrides conflicting guidance in repository templates, agent
instructions, and supporting docs. Amendments MUST update this file and all
affected dependent artifacts in the same change, including a refreshed Sync
Impact Report at the top of the constitution. The derived checklist at
`.specify/memory/constitution-checks.json` is the operational source for plan,
spec, task, and command-template compliance and MUST be updated in the same
change whenever this constitution changes.

Versioning follows semantic rules for governance: MAJOR for removing or
redefining principles in incompatible ways, MINOR for new principles or
materially expanded guidance, and PATCH for clarifications that do not change
expected behavior. Compliance review MUST happen in every plan, spec, task list,
and code review touching this project, and unresolved constitution violations
MUST be documented explicitly and approved before implementation proceeds.

**Version**: 1.0.1 | **Ratified**: 2026-03-07 | **Last Amended**: 2026-03-07
