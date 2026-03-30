---
description: 'Task list for Bitwarden write-back with ds push'
---

# Tasks: Bitwarden write-back with `ds push`

**Input**: Design documents from `/specs/003-ds-push/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are REQUIRED. Load the applicable checks from
`.specify/memory/constitution-checks.json` and include the failing automated
coverage needed to satisfy them.

**Constitution Coverage**: Setup, foundational, story, and polish phases MUST
collectively cover every applicable check ID from
`.specify/memory/constitution-checks.json`.

**Organization**: Tasks are grouped by user story to enable independent
implementation, testing, and constitution compliance validation for each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this belongs to (e.g. `[US1]`, `[US2]`, `[US3]`)
- Include exact file paths in every task description

## Path Conventions

- **CLI entrypoint**: `cmd/ds/`
- **CLI wiring**: `internal/cli/`
- **Config and reporting**: `internal/config/`, `internal/report/`
- **Provider integration**: `internal/provider/bitwarden/`
- **Push and note-json orchestration**: `internal/sync/`
- **Tests and fixtures**: `test/contract/`, `test/integration/`, `test/testdata/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish shared fixtures, helper extensions, and stub support for
push and note-backed provider behavior.

- [x] T001 Create push fixture and golden-output files in `test/testdata/golden/push-dry-run.txt`, `test/testdata/golden/push-fields-mode-error.txt`, `test/testdata/golden/push-unchanged.txt`, and `test/testdata/provider/note-json-valid.json`
- [x] T002 [P] Extend shared CLI helper coverage for push command scenarios in `test/contract/helpers_test.go` and `test/integration/helpers_test.go`
- [x] T003 [P] Extend the `rbw` integration stub to capture note reads and scripted add/edit writes in `test/integration/rbw_stub_test.go`
- [x] T004 [P] Seed malformed note-payload and multiline env fixtures in `test/testdata/provider/note-json-malformed.txt` and `test/testdata/env/push-multiline.env`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build the shared storage-mode, payload, provider-write, and
reporting primitives required by all user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T005 [P] Add failing config and payload unit coverage for `storage_mode` parsing and note-json canonicalization in `internal/config/config_test.go` and `internal/sync/note_json_test.go`
- [x] T006 [P] Add failing provider write-path unit coverage for scripted `rbw` add/edit automation and note reads in `internal/provider/bitwarden/adapter_test.go` and `internal/provider/bitwarden/rbw_client_test.go`
- [x] T007 Implement `storage_mode` defaults and validation in `internal/config/config.go`
- [x] T008 Implement canonical note-json payload types, parsing, and semantic comparison helpers in `internal/sync/note_json.go` and `internal/sync/note_json_test.go`
- [x] T009 Implement write-capable provider contracts and shared push-planning scaffolding in `internal/provider/provider.go`, `internal/sync/engine.go`, and `internal/sync/push.go`
- [x] T010 Implement scriptable `rbw` note read/write primitives in `internal/provider/bitwarden/rbw_client.go`, `internal/provider/bitwarden/adapter.go`, and `internal/provider/bitwarden/rbw_client_test.go`

**Checkpoint**: Storage-mode parsing, canonical payload handling, and provider
write primitives are ready for story work.

---

## Phase 3: User Story 1 - Preview and upload the current `.env` into Bitwarden (Priority: P1) 🎯 MVP

**Goal**: Let users preview and upload the current `.env` into a repo-scoped
Bitwarden item through `ds push` without leaking secrets.

**Independent Test**: Configure a repo for `note_json` storage, stub `rbw`, run
`ds push --dry-run` and `ds push`, then verify preview output is redacted and
the repo-scoped Bitwarden item is created, updated, or left unchanged
correctly.

### Tests for User Story 1 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T011 [P] [US1] Add CLI contract coverage for `ds push` dry-run, create, and unchanged flows in `test/contract/push_contract_test.go`
- [x] T012 [P] [US1] Add integration coverage for `ds push --dry-run`, create, update, and no-op runs in `test/integration/push_command_test.go`
- [x] T013 [P] [US1] Add unit coverage for push planning and command wiring in `internal/sync/push_test.go` and `internal/cli/push_test.go`

### Implementation for User Story 1

- [x] T014 [US1] Implement push diff planning, change classification, and no-op detection in `internal/sync/push.go` and `internal/sync/push_test.go`
- [x] T015 [US1] Implement redacted push summary and per-key output markers in `internal/report/output.go` and `internal/report/redact.go`
- [x] T016 [US1] Implement the `ds push` Cobra command and register it in `internal/cli/push.go` and `internal/cli/root.go`
- [x] T017 [US1] Wire Bitwarden create/update execution and final cache sync into the push flow in `internal/provider/bitwarden/adapter.go`, `internal/provider/bitwarden/rbw_client.go`, and `internal/cli/push.go`
- [x] T018 [US1] Verify push dry-run fidelity and no-op latency budgets in `internal/sync/push_benchmark_test.go` and `test/testdata/golden/push-dry-run.txt`

**Checkpoint**: User Story 1 is independently functional and demoable as the
MVP slice.

---

## Phase 4: User Story 2 - Rebuild and validate from pushed provider data (Priority: P2)

**Goal**: Let the existing read commands consume the note-backed provider data
written by `ds push`.

**Independent Test**: Push a repo-scoped note-json payload through the stubbed
provider, remove `.env`, and verify `ds sync`, `ds diff`, `ds validate`, and
`ds missing` all interpret the same provider state correctly.

### Tests for User Story 2 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T019 [P] [US2] Add contract coverage for note-json-backed `sync`, `diff`, `validate`, and `missing` flows in `test/contract/note_json_contract_test.go` and `test/contract/diff_validate_contract_test.go`
- [x] T020 [P] [US2] Add integration coverage for round-trip reconstruction and malformed payload failures in `test/integration/note_json_sync_test.go` and `test/integration/diff_validate_command_test.go`
- [x] T021 [P] [US2] Add unit coverage for note-json resolution and malformed payload issues in `internal/provider/bitwarden/adapter_test.go`, `internal/sync/forward_test.go`, and `internal/sync/validate_test.go`

### Implementation for User Story 2

- [x] T022 [US2] Implement note-json resolution for `Resolve` and `ResolveMany` in `internal/provider/bitwarden/adapter.go`
- [x] T023 [US2] Teach forward sync planning to read provider-managed values from note-json payloads in `internal/sync/forward.go` and `internal/sync/engine.go`
- [x] T024 [US2] Extend diff, validate, and missing flows for note-json payload issues and redacted previews in `internal/sync/diff.go`, `internal/sync/validate.go`, and `internal/sync/missing.go`
- [x] T025 [US2] Surface note-json payload failures through command handling in `internal/cli/sync.go`, `internal/cli/diff.go`, `internal/cli/validate.go`, and `internal/cli/missing.go`

**Checkpoint**: User Story 2 keeps the pushed provider payload consumable by the
existing read commands.

---

## Phase 5: User Story 3 - Adopt push mode without breaking existing repos (Priority: P3)

**Goal**: Let maintainers opt into push mode deliberately while preserving
backward-compatible field-based reads for current repositories.

**Independent Test**: Run `ds push` in a repo that still uses the default
field-based mode and in a repo with `.env` keys outside `.env.example`, then
verify the CLI explains the mode mismatch, surfaces extras, and leaves the
schema untouched.

### Tests for User Story 3 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T026 [P] [US3] Add contract coverage for fields-mode rejection and extra-key push previews in `test/contract/push_contract_test.go`
- [x] T027 [P] [US3] Add integration coverage for fields-mode rejection and schema-extra previews in `test/integration/push_command_test.go`
- [x] T028 [P] [US3] Add unit coverage for storage-mode defaults and migration guidance in `internal/config/config_test.go` and `internal/cli/push_test.go`

### Implementation for User Story 3

- [x] T029 [US3] Implement fields-mode rejection and actionable migration guidance in `internal/cli/push.go`, `internal/provider/bitwarden/adapter.go`, and `internal/report/output.go`
- [x] T030 [US3] Preserve default field-based read behavior while gating push to `note_json` in `internal/config/config.go` and `internal/provider/bitwarden/adapter.go`
- [x] T031 [US3] Document `storage_mode` opt-in, `ds push`, and migration guidance in `README.md` and `specs/003-ds-push/quickstart.md`
- [x] T032 [US3] Align push help text and config examples with shared CLI wording in `internal/cli/push.go`, `README.md`, and `specs/003-ds-push/contracts/cli-push.contract.md`

**Checkpoint**: User Story 3 makes `ds push` adoptable without silently changing
current field-based repositories.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Finish regression coverage, documentation, redaction audits, and
budget validation across all stories.

- [x] T033 [P] Update push golden outputs and provider fixtures in `test/testdata/golden/push-dry-run.txt`, `test/testdata/golden/push-fields-mode-error.txt`, `test/testdata/golden/push-unchanged.txt`, and `test/testdata/provider/note-json-valid.json`
- [x] T034 [P] Audit secret redaction and shared status vocabulary across push and note-json read paths in `internal/report/output.go`, `internal/report/redact.go`, and `specs/003-ds-push/contracts/cli-push.contract.md`
- [x] T035 Add regression coverage for multiline values, malformed note payloads, and create/update no-op stability in `internal/sync/note_json_test.go`, `test/contract/note_json_contract_test.go`, and `test/integration/push_command_test.go`
- [x] T036 Verify performance and provider-call budgets in `internal/sync/push_benchmark_test.go`, `internal/sync/forward_benchmark_test.go`, and `internal/provider/bitwarden/adapter_test.go`
- [x] T037 Run quickstart validation and refresh operator documentation in `specs/003-ds-push/quickstart.md` and `README.md`

---

## Constitution Coverage

- **workflow-preservation**: T005-T010, T014-T017, T022-T025, T029-T032
- **module-boundaries**: T005-T010, T014, T016-T017, T022-T025, T029-T030
- **deterministic-file-fidelity**: T001-T004, T005, T011-T018, T019-T025, T026-T027, T033-T037
- **test-first-reliability**: T005-T006, T011-T013, T019-T021, T026-T028, T035
- **ux-consistency-secret-safety**: T011-T017, T019-T025, T026-T034
- **performance-dependency-budget**: T005-T010, T018, T024, T036

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion - recommended MVP slice
- **User Story 2 (Phase 4)**: Depends on Foundational completion and reuses the note-json payload contract established for `ds push`
- **User Story 3 (Phase 5)**: Depends on Foundational completion and refines the push-mode adoption path introduced in User Story 1
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: First delivery target and recommended MVP
- **User Story 2 (P2)**: Reuses the payload shape introduced by US1 but remains independently testable once note-json fixtures exist
- **User Story 3 (P3)**: Reuses the push command surface from US1 while keeping existing field-based repos backward compatible

### Within Each User Story

- Tests MUST be written and fail before implementation
- Shared push or note-json helpers precede command or read-path wiring
- Provider primitives precede command execution
- CLI or sync wiring precedes UX/performance validation
- Each story must satisfy its own independent test before moving on

### Parallel Opportunities

- T002-T004 can run in parallel after T001
- T005-T006 can run in parallel before T007-T010
- Within **US1**, T011-T013 can run in parallel before T014-T018
- Within **US2**, T019-T021 can run in parallel before T022-T025
- Within **US3**, T026-T028 can run in parallel before T029-T032
- T033-T034 can run in parallel during polish once the stories are complete

---

## Parallel Example: User Story 1

```bash
# Launch the required tests for User Story 1 together:
Task: "T011 [US1] Add CLI contract coverage in test/contract/push_contract_test.go"
Task: "T012 [US1] Add integration coverage in test/integration/push_command_test.go"
Task: "T013 [US1] Add unit coverage in internal/sync/push_test.go and internal/cli/push_test.go"
```

## Parallel Example: User Story 2

```bash
# Launch the required tests for User Story 2 together:
Task: "T019 [US2] Add note-json contract coverage in test/contract/note_json_contract_test.go and test/contract/diff_validate_contract_test.go"
Task: "T020 [US2] Add note-json integration coverage in test/integration/note_json_sync_test.go and test/integration/diff_validate_command_test.go"
Task: "T021 [US2] Add note-json unit coverage in internal/provider/bitwarden/adapter_test.go, internal/sync/forward_test.go, and internal/sync/validate_test.go"
```

## Parallel Example: User Story 3

```bash
# Launch the required tests for User Story 3 together:
Task: "T026 [US3] Add fields-mode contract coverage in test/contract/push_contract_test.go"
Task: "T027 [US3] Add fields-mode integration coverage in test/integration/push_command_test.go"
Task: "T028 [US3] Add storage-mode guidance unit coverage in internal/config/config_test.go and internal/cli/push_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Confirm `ds push --dry-run` and `ds push` work with a
   note-json repo, stay redacted, and no-op correctly

### Incremental Delivery

1. Finish Setup + Foundational
2. Deliver User Story 1 for immediate push/write-back value
3. Add User Story 2 for note-json-backed round-trip reads
4. Add User Story 3 for safe adoption and backward-compatible guidance
5. Finish polish and regression validation

### Parallel Team Strategy

1. One developer completes Setup + Foundational
2. After Foundational:
   - Developer A owns User Story 1 push previews and writes
   - Developer B owns User Story 2 note-json read compatibility
   - Developer C owns User Story 3 adoption guidance and backward-compatibility validation
3. Rejoin for polish, regression coverage, and quickstart validation

---

## Notes

- [P] tasks use different files or can proceed after shared prerequisites complete
- User story phases map directly to independently testable feature slices
- Suggested MVP scope: **User Story 1 only**
