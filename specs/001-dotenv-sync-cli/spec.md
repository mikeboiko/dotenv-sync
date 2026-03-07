# Feature Specification: dotenv-sync CLI MVP

**Feature Branch**: `001-dotenv-sync-cli`
**Created**: 2026-03-07
**Status**: Draft
**Input**: User description: "Build a CLI application that synchronizes a project's `.env` from `.env.example` using a password manager as the secret source, initially Bitwarden, while preserving the standard developer workflow and providing sync, diff, validate, doctor, init, and schema-maintenance flows."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Generate Local Env Files (Priority: P1)

As a developer joining or returning to a project, I want to create or refresh my
local `.env` file from the committed `.env.example` schema and the team's secret
source in one command so I can start working without changing my normal
workflow.

**Why this priority**: This is the core product promise and the minimum outcome
that makes the tool valuable on day one.

**Independent Test**: In a project with a valid `.env.example`, accessible
Bitwarden secrets, and no existing `.env`, a user can run the sync flow once and
receive a usable `.env` file without wrapping any later application commands.

**Acceptance Scenarios**:

1. **Given** a project contains `.env.example` entries with a mix of blank secret
   slots and explicit static defaults, **When** the user runs the sync flow,
   **Then** the tool creates or updates `.env` by resolving blanks from the
   secret source and copying explicit defaults as written.
2. **Given** `.env` already matches the schema and resolved secrets, **When**
   the user runs the sync flow again, **Then** the tool reports that the file is
   already up to date and leaves the file unchanged.
3. **Given** the secret source is unavailable, locked, or missing a required
   value, **When** the user runs the sync flow, **Then** the tool stops safely,
   explains what is wrong, and tells the user how to recover without revealing
   any secret values.

---

### User Story 2 - Preview and Validate Changes (Priority: P2)

As a developer or team maintainer, I want to preview and validate differences
between `.env.example`, `.env`, and the secret source so I can trust changes,
catch drift early, and use the tool in CI or onboarding workflows.

**Why this priority**: Teams need to trust file changes before writing them and
need a way to detect drift or missing prerequisites before runtime failures.

**Independent Test**: With intentionally mismatched `.env.example`, `.env`, and
secret availability, a user can run preview, diff, or validate flows and see
accurate, redacted results without writing files.

**Acceptance Scenarios**:

1. **Given** the sync flow would change one or more variables, **When** the user
   runs a dry-run or diff flow, **Then** the tool shows only the real additions,
   updates, unchanged entries, and missing values, with secret content redacted.
2. **Given** `.env` or `.env.example` has missing, extra, or malformed entries,
   **When** the user runs validate, **Then** the tool reports the mismatch,
   explains the impact, and exits in a way that local automation and CI can use.
3. **Given** a required secret cannot be resolved from the provider, **When**
   the user runs a missing-secrets or validation flow, **Then** the tool lists
   the unresolved keys without exposing any secret material.

---

### User Story 3 - Bootstrap and Maintain the Schema (Priority: P3)

As a maintainer, I want to initialize a schema from an existing local file,
diagnose provider prerequisites, and optionally backfill new keys into the
schema so the project stays consistent as configuration evolves.

**Why this priority**: This extends the tool from one-time sync into an ongoing
team maintenance workflow without undermining the schema contract.

**Independent Test**: Starting from a project with only `.env`, or from a
project where `.env` contains keys missing from `.env.example`, a maintainer can
run init, doctor, and opt-in schema-maintenance flows and get a safe schema or
clear diagnostic result.

**Acceptance Scenarios**:

1. **Given** a project has a populated `.env` but no `.env.example`, **When**
   the maintainer runs init, **Then** the tool generates a schema file that
   keeps safe explicit defaults and strips secret values to blank placeholders.
2. **Given** the provider CLI is missing, the user is logged out, or the vault
   is locked, **When** the maintainer runs doctor, **Then** the tool identifies
   each failed prerequisite and the next action required to fix it.
3. **Given** `.env` contains additional keys that are not yet in
   `.env.example`, **When** the maintainer runs the opt-in reverse-sync flow,
   **Then** the tool adds those keys to `.env.example` as blank placeholders
   without copying secret values into the schema.

### Edge Cases

- How does the system behave when `.env.example` contains duplicate keys,
  conflicting comments, or malformed lines?
- What happens when secret values contain spaces, quotes, multiline content, or
  characters that require careful file formatting?
- How does the tool respond if the provider session expires after some keys have
  been resolved but before the run completes?
- What happens when `.env` already contains comments or manual ordering that do
  not match the schema exactly?
- How are unmapped variables handled when the provider naming convention differs
  from the schema key names?
- What happens when the user requests reverse sync but `.env.example` contains
  hand-written explanatory comments that must remain intact?

## Assumptions

- The first release targets local developer and CI workflows rather than secret
  rotation or write-back to the provider.
- `.env.example` is the canonical schema for sync, diff, and validate unless the
  user is explicitly running init to create that schema.
- Blank values in `.env.example` represent provider-managed secrets, while
  explicit values represent safe defaults that may be copied as written.
- Reverse sync is opt-in and is limited to adding blank placeholders for new
  keys discovered in `.env`.
- Bitwarden is the initial provider available out of the box, and future
  providers will be added without changing the standard `.env` workflow.

## User Experience Consistency _(mandatory)_

- **UX-001**: Command output MUST use a shared status vocabulary that clearly
  distinguishes detected prerequisites, resolved values, unchanged entries,
  missing values, written updates, and failures.
- **UX-002**: Errors MUST explain the failed prerequisite or mismatch, why it
  blocks the requested action, and the next recovery action the user should take.
- **UX-003**: Secret-bearing values MUST be redacted or represented only by
  their resolution status in logs, previews, diffs, diagnostics, and errors.
- **UX-004**: Dry-run and diff output MUST show only real mutations and MUST
  make it obvious which entries are added, updated, unchanged, or unresolved.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST treat `.env.example` as the canonical schema for
  expected environment variables during sync, diff, and validate flows.
- **FR-002**: The system MUST create or update `.env` by copying explicit
  defaults from `.env.example` and resolving blank schema entries from the
  configured secret source.
- **FR-003**: Users MUST be able to produce a usable local `.env` without
  wrapping their normal application commands or changing how tools already load
  `.env` files.
- **FR-004**: The system MUST support Bitwarden as the initial secret source and
  MUST verify prerequisite states required to retrieve secrets.
- **FR-005**: The system MUST provide a dry-run preview before writing changes
  so users can understand what will happen in advance.
- **FR-006**: The system MUST provide diff and validate flows that identify
  missing variables, extra variables, unresolved secrets, and schema drift
  across `.env.example`, `.env`, and the secret source.
- **FR-007**: The system MUST preserve comments, key ordering, and unchanged
  lines when rewriting `.env` or `.env.example` wherever the source format
  allows.
- **FR-008**: The system MUST provide a diagnostic flow that reports missing
  provider CLI installation, logged-out sessions, locked vaults, and other
  prerequisite failures with actionable recovery guidance.
- **FR-009**: The system MUST support configurable key mapping so a schema key
  can resolve from a differently named provider record when teams need it.
- **FR-010**: The system MUST provide an opt-in reverse-sync flow that adds new
  keys from `.env` back into `.env.example` as blank placeholders only.
- **FR-011**: The system MUST provide an init flow that can generate
  `.env.example` from an existing `.env` while stripping secret values and
  preserving safe static defaults.
- **FR-012**: The system MUST provide a missing-secrets report that lists schema
  keys that cannot currently be resolved from the provider without exposing their
  values.
- **FR-013**: The system MUST produce machine-usable success and failure
  outcomes so local automation and CI can rely on validate and diagnostic flows.
- **FR-014**: The system MUST keep schema-maintenance flows from writing secret
  values into `.env.example`.

### Key Entities _(include if feature involves data)_

- **Environment Schema**: The committed description of expected environment
  variables, including variable names, blank secret placeholders, safe default
  literals, comments, and intended ordering.
- **Local Environment File**: The developer-local `.env` file that contains the
  resolved values and defaults consumed by frameworks, tools, and editors.
- **Secret Resolution Record**: The outcome of attempting to satisfy a schema
  entry, including whether the value came from a static default, a provider
  mapping, or remains unresolved.
- **Provider Readiness State**: The current availability of the configured
  secret source, including whether required tooling is installed, authenticated,
  and unlocked.

## Performance Requirements _(mandatory)_

- **PR-001**: Users MUST receive local sync, diff, validate, or dry-run results
  for projects with up to 500 variables within 2 seconds when provider access is
  already available.
- **PR-002**: Users MUST receive provider-backed sync or validation results for
  projects with up to 500 variables within 10 seconds under normal provider
  availability, and a single run MUST not require more than one provider lookup
  per distinct variable unless reauthentication or pagination is required.
- **PR-003**: No-op sync runs for projects with up to 500 variables MUST finish
  within 1 second and MUST not modify file contents when no semantic changes are
  needed.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: A new developer can generate a usable `.env` for a project that
  already has `.env.example` and accessible secrets in one command within
  5 minutes of entering the project directory.
- **SC-002**: At least 95% of repeated no-op sync runs on representative
  projects report that the local environment file is up to date without changing
  file contents.
- **SC-003**: Acceptance testing finds zero raw secret exposures across preview,
  diff, validation, diagnostic, and error output.
- **SC-004**: At least 90% of seeded schema-drift, missing-secret, or missing-
  prerequisite scenarios are identified before application startup through the
  tool's preview, validation, missing, or diagnostic flows.
