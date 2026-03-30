# Implementation Plan: Bitwarden write-back with `ds push`

**Branch**: `003-ds-push` | **Date**: 2026-03-30 | **Spec**:
[spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-ds-push/spec.md`

## Summary

Add an opt-in `ds push` command that uploads the current `.env` into a
repo-scoped Bitwarden entry without printing secret values, while keeping the
existing field-based Bitwarden layout untouched for current repositories. The
design uses a new `note_json` storage mode backed by the repo item's Bitwarden
notes because `rbw` already supports scriptable note/password edits and note
reads, whereas the current field-based `rbw` flow is read-oriented. Existing
read commands (`sync`, `diff`, `validate`, `missing`) will gain compatibility
with the note-backed payload when that mode is configured.

## Technical Context

**Language/Version**: Go 1.22  
**Primary Dependencies**: `github.com/spf13/cobra` for CLI routing, `gopkg.in/yaml.v3` for config loading, Go standard-library JSON support, installed `rbw` CLI  
**Storage**: Local env files plus a repo-scoped Bitwarden login-item notes payload in `note_json` mode  
**Testing**: `go test ./...`, CLI contract tests, subprocess integration tests
with `rbw` stubs, unit tests for payload serialization/diffing and editor-helper
flows  
**Target Platform**: Cross-platform CLI for Linux, macOS, and Windows on
developer machines with `rbw` installed  
**Project Type**: Single-binary CLI with provider-backed local workflow support  
**Performance Goals**: `ds push --dry-run` completes within 200 ms p95 for
`.env` files up to 500 keys with stubbed provider responses; real `ds push`
uses at most one provider read, one add-or-edit mutation, and one final sync  
**Constraints**: No runtime wrappers, no secret values in output, no mutation of
`.env` or `.env.example` during push, backward compatibility for existing
field-based repos, no new heavy runtime dependencies, and `rbw` write behavior
must remain scriptable through its editor-based CLI surface  
**Scale/Scope**: One repo-scoped Bitwarden item per repository, up to roughly
500 env keys per repo payload, one new push command, one new storage mode, and
read compatibility for the existing provider-backed commands

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

Pre-design review passes because the feature adds an opt-in provider write-back
path without changing the existing `.env` workflow or the default Bitwarden
layout. Post-design review also passes after `research.md`, `data-model.md`,
`quickstart.md`, and the contracts below define a testable, redacted, and
backward-compatible path for push plus note-backed reads.

| Check ID                      | Status | Evidence                                                                                                                                                                                                        |
| ----------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| workflow-preservation         | PASS   | The Summary and Technical Context keep `.env` as the developer-facing file, preserve `.env.example` as schema context, and scope write-back to an explicit `ds push` command with no runtime wrapper changes.   |
| module-boundaries             | PASS   | The Project Structure keeps CLI routing in `internal/cli`, config in `internal/config`, write/read orchestration in `internal/sync`, and Bitwarden-specific mutation behavior in `internal/provider/bitwarden`. |
| deterministic-file-fidelity   | PASS   | `ds push` does not rewrite local env files, and `research.md` Decisions 2 and 3 plus `contracts/cli-push.contract.md` require deterministic previews and no-op stability for provider payload changes.          |
| test-first-reliability        | PASS   | Technical Context requires unit, contract, and integration tests, while `research.md` Decisions 4 and 5 identify failing-first coverage for editor automation, malformed payloads, and provider failures.       |
| ux-consistency-secret-safety  | PASS   | `spec.md` UX requirements and `contracts/cli-push.contract.md` enforce shared status vocabulary, actionable errors, and fully redacted output across preview and write paths.                                   |
| performance-dependency-budget | PASS   | Technical Context and `research.md` Decision 6 constrain the provider call budget to one read, one mutation, and one sync while avoiding new heavy runtime dependencies.                                        |

## Project Structure

### Documentation (this feature)

```text
specs/003-ds-push/
├── plan.md
├── spec.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cli-push.contract.md
│   └── provider-writeback.contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── ds/
    └── main.go

internal/
├── cli/
│   ├── root.go
│   └── push.go
├── config/
│   └── config.go
├── provider/
│   ├── provider.go
│   └── bitwarden/
│       ├── adapter.go
│       ├── rbw_client.go
│       └── status.go
├── report/
│   ├── output.go
│   └── redact.go
└── sync/
    ├── engine.go
    ├── push.go
    └── note_json.go

test/
├── contract/
│   ├── helpers_test.go
│   ├── push_contract_test.go
│   └── note_json_contract_test.go
├── integration/
│   ├── helpers_test.go
│   ├── rbw_stub_test.go
│   ├── push_command_test.go
│   └── note_json_sync_test.go
└── testdata/
    └── golden/

README.md
```

**Structure Decision**: Extend the existing single-binary CLI structure by
adding one focused command handler, one push-planning path in `internal/sync`,
and one write-capable Bitwarden adapter path. This keeps provider mutation logic
out of Cobra handlers, lets existing read commands share the same storage-mode
awareness, and preserves the repo's current test split between contract and
integration coverage.

## Complexity Tracking

No constitution violations or complexity exceptions are expected for this
feature.
