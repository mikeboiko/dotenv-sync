---
description: 'Task list for implementing dotenv-sync CLI MVP'
---

# Tasks: dotenv-sync CLI MVP

**Input**: Design documents from `/specs/001-dotenv-sync-cli/`
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
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **CLI entrypoint**: `cmd/ds/`
- **Application packages**: `internal/cli/`, `internal/config/`, `internal/envfile/`, `internal/fs/`, `internal/provider/`, `internal/report/`, `internal/sync/`
- **Public package**: `pkg/dotenvsync/`
- **Tests and fixtures**: `test/contract/`, `test/integration/`, `test/testdata/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize the Go module, the `ds` executable entrypoint, and the
shared test/fixture scaffolding used by every story.

- [x] T001 Initialize the Go module and default `ds` entrypoint in `go.mod`, `go.sum`, `cmd/ds/main.go`, and `pkg/dotenvsync/version.go`
- [x] T002 [P] Scaffold the Cobra root and subcommand files in `internal/cli/root.go`, `internal/cli/sync.go`, `internal/cli/diff.go`, `internal/cli/validate.go`, `internal/cli/doctor.go`, `internal/cli/init.go`, `internal/cli/missing.go`, and `internal/cli/reverse.go`
- [x] T003 [P] Create shared Go test helpers in `test/contract/helpers_test.go`, `test/integration/helpers_test.go`, and `test/testdata/README.md`
- [x] T004 [P] Seed baseline env, provider, and golden fixtures in `test/testdata/env/basic.env.example`, `test/testdata/env/basic.env`, `test/testdata/provider/rbw-list-success.txt`, and `test/testdata/golden/sync-dry-run.txt`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build the shared infrastructure that every command relies on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T005 Implement `.envsync.yaml` loading and default resolution in `internal/config/config.go`
- [x] T006 [P] Implement shared output vocabulary, redaction, and exit-code helpers in `internal/report/output.go`, `internal/report/redact.go`, and `internal/report/exitcodes.go`
- [x] T007 [P] Implement envfile domain types and token parsing in `internal/envfile/model.go` and `internal/envfile/parser.go`
- [x] T008 [P] Implement deterministic writing, merge helpers, and atomic rewrites in `internal/envfile/writer.go`, `internal/envfile/merge.go`, and `internal/fs/atomic.go`
- [x] T009 [P] Define the provider interface and base `rbw` client primitives in `internal/provider/provider.go`, `internal/provider/bitwarden/rbw_client.go`, and `internal/provider/bitwarden/status.go`
- [x] T010 Implement shared sync planning and change classification in `internal/sync/engine.go`
- [x] T011 [P] Add reusable `rbw` stub wiring for integration tests in `test/integration/rbw_stub_test.go`, `test/testdata/provider/rbw-status-unlocked.txt`, and `test/testdata/provider/rbw-get-database-url.txt`
- [x] T012 [P] Add benchmark helpers and 500-key fixtures in `test/integration/benchmark_helpers_test.go` and `test/testdata/env/large-schema.env.example`

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Generate Local Env Files (Priority: P1) 🎯 MVP

**Goal**: Let a developer create or refresh `.env` from `.env.example` and
Bitwarden via `rbw` without changing the standard workflow.

**Independent Test**: Run `ds sync` in a fixture project with blank secret
placeholders, static defaults, and an `rbw` stub; verify that `.env` is written,
static defaults are copied, provider-managed keys resolve, and a second run is a
true no-op.

### Tests for User Story 1 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T013 [P] [US1] Add CLI contract coverage for `ds sync` success, no-op, and `--dry-run` output in `test/contract/sync_contract_test.go`
- [x] T014 [P] [US1] Add integration coverage for forward sync, static defaults, and no-op reruns in `test/integration/sync_command_test.go`
- [x] T015 [P] [US1] Add unit coverage for forward planning and deterministic writer stability in `internal/sync/forward_test.go` and `internal/envfile/writer_test.go`

### Implementation for User Story 1

- [x] T016 [P] [US1] Implement the `rbw`-backed Bitwarden resolution adapter and per-run caching in `internal/provider/bitwarden/adapter.go`
- [x] T017 [US1] Implement forward-sync planning and apply logic in `internal/sync/forward.go`
- [x] T018 [US1] Implement the `ds sync` command path in `internal/cli/sync.go`
- [x] T019 [US1] Wire sync-specific status summaries and redacted rendering in `internal/report/output.go` and `internal/report/redact.go`
- [x] T020 [US1] Add sync performance benchmarks and lookup-budget assertions in `internal/sync/forward_benchmark_test.go` and `test/integration/sync_benchmark_test.go`

**Checkpoint**: User Story 1 is independently functional, testable, and ready to demo.

---

## Phase 4: User Story 2 - Preview and Validate Changes (Priority: P2)

**Goal**: Let users preview, diff, validate, and report unresolved values
without writing files.

**Independent Test**: Run `ds diff`, `ds validate`, and `ds missing` against
fixtures with drift, malformed input, and missing secrets; verify redacted,
CI-friendly results with correct exit codes.

### Tests for User Story 2 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T021 [P] [US2] Add CLI contract coverage for `ds diff`, `ds validate`, and `ds missing` in `test/contract/diff_validate_contract_test.go`
- [x] T022 [P] [US2] Add integration coverage for drift, duplicate keys, malformed files, and unresolved secrets in `test/integration/diff_validate_command_test.go`
- [x] T023 [P] [US2] Add unit coverage for diff, validation, and missing-key classification in `internal/sync/diff_test.go`, `internal/sync/validate_test.go`, and `internal/sync/missing_test.go`

### Implementation for User Story 2

- [x] T024 [US2] Implement diff, validate, and missing planners in `internal/sync/diff.go`, `internal/sync/validate.go`, and `internal/sync/missing.go`
- [x] T025 [US2] Implement `ds diff`, `ds validate`, and `ds missing` command handlers in `internal/cli/diff.go`, `internal/cli/validate.go`, and `internal/cli/missing.go`
- [x] T026 [US2] Extend CI-friendly exit handling and preview formatting in `internal/report/output.go` and `internal/report/exitcodes.go`
- [x] T027 [US2] Add diff and validate benchmark coverage for 500-key fixtures in `internal/sync/diff_benchmark_test.go` and `internal/sync/validate_benchmark_test.go`

**Checkpoint**: User Stories 1 and 2 both work independently, with preview and validation flows ready for local use and CI.

---

## Phase 5: User Story 3 - Bootstrap and Maintain the Schema (Priority: P3)

**Goal**: Let maintainers bootstrap `.env.example`, diagnose `rbw`
prerequisites, and add new schema placeholders safely.

**Independent Test**: Run `ds doctor`, `ds init`, and `ds reverse` against
fixtures with missing schema files, locked or logged-out `rbw` states, and extra
local env keys; verify actionable diagnostics and blank-placeholder schema
writes only.

### Tests for User Story 3 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T028 [P] [US3] Add CLI contract coverage for `ds doctor`, `ds init`, and `ds reverse` in `test/contract/doctor_init_reverse_contract_test.go`
- [x] T029 [P] [US3] Add integration coverage for init schema generation, reverse-sync placeholders, and `rbw` prerequisite failures in `test/integration/doctor_init_reverse_test.go`
- [x] T030 [P] [US3] Add unit coverage for provider readiness and schema-maintenance planning in `internal/provider/bitwarden/status_test.go` and `internal/envfile/merge_test.go`

### Implementation for User Story 3

- [x] T031 [US3] Implement `rbw` readiness checks and diagnostic mapping in `internal/provider/bitwarden/status.go`
- [x] T032 [US3] Implement init and reverse-sync planning in `internal/sync/reverse.go` and `internal/envfile/merge.go`
- [x] T033 [US3] Implement `ds doctor`, `ds init`, and `ds reverse` command handlers in `internal/cli/doctor.go`, `internal/cli/init.go`, and `internal/cli/reverse.go`
- [x] T034 [US3] Validate no-secret schema writes and operator guidance with golden outputs in `test/testdata/golden/doctor-locked.txt`, `test/testdata/golden/init-preview.txt`, and `internal/report/output.go`

**Checkpoint**: All three user stories are independently functional and maintain the schema contract safely.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Complete documentation, regression coverage, and cross-story validation.

- [x] T035 [P] Update user-facing command documentation in `README.md` and `specs/001-dotenv-sync-cli/quickstart.md`
- [x] T036 [P] Run and tune regression benchmarks in `internal/sync/forward_benchmark_test.go`, `internal/sync/diff_benchmark_test.go`, and `internal/sync/validate_benchmark_test.go`
- [x] T037 [P] Audit redaction and error-code consistency in `internal/report/redact.go`, `internal/report/output.go`, `test/contract/sync_contract_test.go`, `test/contract/diff_validate_contract_test.go`, and `test/contract/doctor_init_reverse_contract_test.go`
- [x] T038 Validate cross-platform path, line-ending, and temp-file behavior in `internal/fs/atomic.go`, `internal/envfile/writer.go`, `test/integration/sync_command_test.go`, and `test/integration/doctor_init_reverse_test.go`
- [x] T039 Run quickstart verification and refresh sample fixtures in `specs/001-dotenv-sync-cli/quickstart.md`, `test/testdata/env/basic.env.example`, `test/testdata/env/basic.env`, and `test/testdata/env/no-schema.env`

---

## Constitution Coverage

- **workflow-preservation**: T005, T010, T013-T018, T031-T033
- **module-boundaries**: T001-T010
- **deterministic-file-fidelity**: T007-T008, T014-T015, T022-T024, T029-T032, T038
- **test-first-reliability**: T013-T015, T021-T023, T028-T030, T036
- **ux-consistency-secret-safety**: T006, T019, T026, T031-T034, T037
- **performance-dependency-budget**: T012, T020, T027, T036, T038

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion - MVP scope
- **User Story 2 (Phase 4)**: Depends on Foundational completion and may reuse User Story 1 sync internals, but remains independently testable through diff/validate flows
- **User Story 3 (Phase 5)**: Depends on Foundational completion and may reuse provider/config foundations from User Story 1, but remains independently testable through doctor/init/reverse flows
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: First deliverable and the recommended MVP slice
- **User Story 2 (P2)**: Builds on the same schema, provider, and report primitives but can be validated without User Story 3
- **User Story 3 (P3)**: Builds on provider readiness and envfile fidelity work but can be validated independently from User Story 2 once the foundations exist

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Provider and planner changes precede command wiring
- Command wiring precedes UX/performance validation
- Story-specific benchmarks and golden checks complete the story before sign-off

### Parallel Opportunities

- T002-T004 can run in parallel after T001 establishes the module
- T006-T009 and T011-T012 can run in parallel once T005 establishes config defaults
- Within **US1**, T013-T015 can run in parallel, then T016 can proceed alongside T017 once the tests exist
- Within **US2**, T021-T023 can run in parallel before T024-T027 sequence through planners, commands, and performance checks
- Within **US3**, T028-T030 can run in parallel before T031-T034 sequence through readiness, schema maintenance, and golden validation
- T035-T039 can be split across team members after the desired stories are complete

---

## Parallel Example: User Story 1

```bash
# Launch the required tests for User Story 1 together:
Task: "T013 [US1] Add CLI contract coverage for ds sync in test/contract/sync_contract_test.go"
Task: "T014 [US1] Add integration coverage for sync in test/integration/sync_command_test.go"
Task: "T015 [US1] Add unit coverage for forward planning in internal/sync/forward_test.go and internal/envfile/writer_test.go"

# After the tests fail, split the implementation by file ownership:
Task: "T016 [US1] Implement internal/provider/bitwarden/adapter.go"
Task: "T017 [US1] Implement internal/sync/forward.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Run `go test ./... -run Sync` and the contract/integration suites for `ds sync`
5. Demo the MVP by generating `.env` from `.env.example` with the `rbw` stub

### Incremental Delivery

1. Complete Setup + Foundational → foundation ready
2. Add User Story 1 → validate forward sync and no-op behavior
3. Add User Story 2 → validate preview, drift detection, and CI exit handling
4. Add User Story 3 → validate doctor, init, and reverse-sync flows
5. Finish polish work → validate docs, redaction, and performance budgets

### Parallel Team Strategy

With multiple developers:

1. Developer A: Setup + config/report foundations
2. Developer B: Envfile parser/writer + atomic file helpers
3. Developer C: Provider interface + `rbw` integration harness
4. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
5. Rejoin for Phase 6 polish and regression validation

---

## Notes

- [P] tasks = different files, no blocking dependency on unfinished work
- Each user story includes failing-first tests, implementation, and UX/performance validation
- The default executable name is `ds`; the product and repository name remain `dotenv-sync`
- `bw` compatibility is intentionally deferred to the roadmap and is not part of the MVP task list
- All tasks reference concrete file paths so an implementation agent can execute them without reopening the plan
