# Feature Specification: Bitwarden write-back with `ds push`

**Feature Branch**: `003-ds-push`  
**Created**: 2026-03-30  
**Status**: Draft  
**Input**: User description: "I want to `cd` into a repository such as `~/git/Jesse` and use `ds` to save the environment variables currently in `.env` to Bitwarden."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Preview and upload the current `.env` into Bitwarden (Priority: P1)

As a developer working inside a project repository, I want `ds push` to preview
and upload the values from my current `.env` file into a repo-scoped Bitwarden
entry so I can back up or share the current local configuration without manual
copy/paste.

**Why this priority**: This is the primary user value. Without a direct push
command, operators still have to copy secrets into Bitwarden by hand.

**Independent Test**: Configure a project for Bitwarden note-backed storage,
stub `rbw`, run `ds push --dry-run` and `ds push`, then verify the repo-scoped
Bitwarden entry is created or updated while command output remains redacted.

**Acceptance Scenarios**:

1. **Given** `.env` exists and the repo is configured for Bitwarden
   `note_json` storage, **When** the user runs `ds push --dry-run`, **Then** the
   CLI previews only the real adds, updates, unchanged entries, and extras that
   would be stored, and exits with code `0` without mutating either local files
   or Bitwarden.
2. **Given** no repo-scoped Bitwarden entry exists yet, **When** the user runs
   `ds push`, **Then** the CLI creates the entry and stores the current `.env`
   values in a deterministic provider payload.
3. **Given** the Bitwarden payload already matches `.env`, **When** the user
   runs `ds push`, **Then** the CLI reports `UNCHANGED` and performs no provider
   mutation.

---

### User Story 2 - Rebuild and validate from pushed provider data (Priority: P2)

As a developer, I want the values uploaded by `ds push` to remain consumable by
`ds sync`, `ds diff`, `ds validate`, and `ds missing` so I can round-trip a
repo between local `.env` files and Bitwarden without inventing a second secret
format.

**Why this priority**: Uploading secrets is only useful if the existing read
commands can later consume what was written.

**Independent Test**: Push a repo-scoped Bitwarden payload through the provider
stub, remove `.env`, and verify `ds sync` reconstructs it while `ds diff`,
`ds validate`, and `ds missing` interpret the same provider state correctly.

**Acceptance Scenarios**:

1. **Given** a repo-scoped Bitwarden payload produced by `ds push`, **When** the
   user runs `ds sync`, **Then** `.env` is rebuilt from `.env.example` plus the
   pushed provider values.
2. **Given** the provider payload has been manually corrupted, **When** the user
   runs `ds sync` or `ds validate`, **Then** the CLI fails with an actionable
   error that explains the malformed payload and the recovery path.
3. **Given** the provider payload differs from `.env`, **When** the user runs
   `ds diff`, **Then** the CLI shows only the real changes and never prints raw
   values.

---

### User Story 3 - Use shared field aliases without breaking existing repos (Priority: P3)

As a maintainer, I want `ds push` to support the existing field-based Bitwarden
layout when a repository uses shared password-field aliases so I can reuse one
Bitwarden value across repos without silently breaking current read behavior.

**Why this priority**: The repository already ships with a field-based default,
so adding write-back should support the shared-alias workflow while still
failing safely for unsupported custom-field writes.

**Independent Test**: Run `ds push` in a repo that uses the default field-based
storage mode with `mapping: password`, then run it in a repo whose mapping
targets a custom field and whose `.env` contains keys outside the schema; verify
the shared alias writes succeed, unsupported mappings fail with guidance, and
schema extras are surfaced without mutating `.env.example`.

**Acceptance Scenarios**:

1. **Given** the repository still uses the default Bitwarden `fields` storage
   mode and its pushed keys map to Bitwarden's built-in `password` field,
   **When** the user runs `ds push`, **Then** the CLI updates that shared field
   without requiring `note_json`.
2. **Given** `.env` contains keys that are not present in `.env.example`,
   **When** the user runs `ds push --dry-run`, **Then** those keys are surfaced
   as `extra` entries in the preview, but `.env.example` is not modified.
3. **Given** the repository stays on `fields` mode but maps pushed keys to
   custom Bitwarden fields, **When** the user runs `ds push`, **Then** the CLI
   fails with actionable guidance to use the `password` field or switch to
   `note_json`.

### Edge Cases

- What happens when `.env` is missing, malformed, or contains duplicate keys?
- What happens when the repo-scoped Bitwarden item already exists but contains a
  malformed note payload or an unrelated manual note?
- How are multiline, quoted, or whitespace-sensitive dotenv values serialized so
  they round-trip correctly through provider storage?
- What happens when `.env` contains keys that are intentionally absent from
  `.env.example`?
- How does `ds push` behave when the vault is locked or the item cannot be
  created because `rbw` cannot launch its editor workflow?

## User Experience Consistency _(mandatory)_

- **UX-001**: `ds push` MUST use the repository's existing status vocabulary
  (`CHECKED`, `WRITTEN`, `UNCHANGED`, `ERROR`) and key-level change wording that
  matches other commands.
- **UX-002**: Push failures MUST explain the blocked prerequisite or incompatible
  storage mode, why it prevents upload, and the next recovery step.
- **UX-003**: Push previews, summaries, and failures MUST never print raw secret
  values; they may reference key names and redacted markers only.
- **UX-004**: `ds push --dry-run` MUST describe exactly the same provider
  mutations that a real `ds push` would perform, excluding the final write.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST provide a `ds push` command with the existing
  global `--config`, `--schema`, and `--env` flags plus a `--dry-run` flag.
- **FR-002**: `ds push` MUST read the current `.env` file as the upload source
  and MUST NOT mutate `.env` or `.env.example`.
- **FR-003**: `ds push` MUST evaluate `.env` against `.env.example`, surface
  extra local keys in previews and summaries, and keep `.env.example` unchanged.
- **FR-004**: Configuration MUST support a Bitwarden `storage_mode` with
  `fields` as the backward-compatible default and `note_json` as the opt-in
  write-back mode.
- **FR-005**: In `note_json` mode, the Bitwarden adapter MUST store the repo's
  current `.env` values in a deterministic JSON payload inside the repo-scoped
  Bitwarden item's notes.
- **FR-006**: `ds push` MUST create the repo-scoped Bitwarden item when it does
  not exist and MUST update it only when the semantic payload changes.
- **FR-007**: `ds sync`, `ds diff`, `ds validate`, and `ds missing` MUST support
  resolving values from the `note_json` payload when that storage mode is
  configured.
- **FR-008**: Existing field-based Bitwarden behavior MUST remain unchanged for
  repositories that do not opt into `note_json`.
- **FR-009**: `ds push` MUST fail with actionable diagnostics when `.env` is
  missing or malformed, when the Bitwarden storage mode is incompatible, when
  the provider payload cannot be parsed, or when Bitwarden is unavailable.
- **FR-010**: Operator-visible output for push and note-backed resolution MUST
  never print secret values.
- **FR-011**: The provider abstraction MUST remain modular so future providers or
  future Bitwarden write strategies can implement the same push contract.
- **FR-012**: README, quickstart, and command help MUST document how to opt into
  `note_json` mode and how to use `ds push` safely.
- **FR-013**: In `fields` mode, `ds push` MUST support updating
  provider-managed keys currently present in `.env` when they map to
  Bitwarden's built-in `password` field.
- **FR-014**: In `fields` mode, `ds push` MUST fail with actionable diagnostics
  when pushed keys map to unsupported custom Bitwarden fields or collapse
  conflicting local values into one shared field.

### Key Entities _(include if feature involves data)_

- **Local Env Snapshot**: The parsed `.env` content that serves as the upload
  source for `ds push`.
- **Push Storage Mode**: The configured Bitwarden storage strategy, either
  backward-compatible `fields` mode or write-capable `note_json` mode.
- **Provider Env Payload**: The repo-scoped Bitwarden representation of the env
  values, including whether the item exists and whether the notes payload is
  valid.
- **Push Plan**: The computed diff between `.env`, `.env.example`, and the
  provider payload, including previewable changes and whether a write is needed.

## Performance Requirements _(mandatory)_

- **PR-001**: `ds push --dry-run` MUST complete within 200 ms p95 for `.env`
  files up to 500 keys when provider responses are stubbed or cached.
- **PR-002**: A real `ds push` in `note_json` mode MUST use no more than one
  provider read, one provider create-or-update mutation, and one final `rbw
sync` per command.
- **PR-003**: The feature MUST not add new runtime dependencies beyond the
  existing Go module footprint and standard-library JSON support.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: A developer can upload the current `.env` for a repo into
  Bitwarden with one `ds push` command after opting into `note_json` mode.
- **SC-002**: A second `ds push` against unchanged local and provider state
  returns `UNCHANGED` with no provider mutation.
- **SC-003**: A repo configured for `note_json` mode can delete `.env`, run
  `ds sync`, and recover the same values from Bitwarden without manual copying.
- **SC-004**: Existing repositories that stay on `fields` mode continue to work
  unchanged and receive explicit guidance instead of silent behavior changes.
