---
description: 'Task list for automatic patch release automation'
---

# Tasks: Automatic patch release automation

**Input**: Design documents from `/specs/002-versioning-release-automation/`
**Prerequisites**: `plan.md` (required), `spec.md` (required for user stories), `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Tests are REQUIRED. Each user story begins with failing automated
coverage and includes negative-path validation for release skips and failures.

**Constitution Coverage**: Setup, foundational, story, and polish phases
collectively cover every applicable check ID from
`.specify/memory/constitution-checks.json`.

**Organization**: Tasks are grouped by user story so each slice can be built,
tested, and validated independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel when work touches different files and does not
  depend on incomplete tasks
- **[Story]**: Required on user-story tasks only (`[US1]`, `[US2]`, `[US3]`)
- Every task includes exact repository file paths

## Path Conventions

- Release workflow: `.github/workflows/`
- Release logic: `internal/release/`
- Local preview helper: `scripts/nextversion/`
- Version verification: `pkg/dotenvsync/`, `cmd/ds/`
- Tests and fixtures: `test/contract/`, `test/integration/`, `test/testdata/`
- Documentation: `README.md`, `specs/002-versioning-release-automation/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish fixtures, helpers, and golden outputs used across
automatic patch release planning and validation.

- [X] T001 Create automatic release fixtures in `test/testdata/release/no-tags.txt`, `test/testdata/release/mixed-tags.txt`, `test/testdata/release/already-released.txt`, and `test/testdata/release/concurrent-main-push.txt`
- [X] T002 [P] Extend release workflow helpers in `test/contract/helpers_test.go` and `test/integration/helpers_test.go`
- [X] T003 [P] Capture golden release outcomes in `test/testdata/golden/release-published.txt`, `test/testdata/golden/release-skipped.txt`, and `test/testdata/release/expected-assets.txt`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build shared patch-calculation and release-state primitives that all
stories depend on.

**⚠️ CRITICAL**: No user story work should begin until this phase is complete.

- [X] T004 [P] Add failing unit coverage for patch-only next-version calculation and non-semver filtering in `internal/release/semver_test.go`
- [X] T005 [P] Add failing unit coverage for already-released commit detection, branch validation, and release-state helpers in `internal/release/semver_test.go`
- [X] T006 Implement patch-only version calculation and release-state helpers in `internal/release/semver.go`
- [X] T007 Implement local next-patch preview defaults and release-logic isolation in `scripts/nextversion/main.go`

**Checkpoint**: Shared patch calculation and release-state helpers are ready for
story work.

---

## Phase 3: User Story 1 - Publish the next patch release on every push to `main` (Priority: P1) 🎯 MVP

**Goal**: Publish a new patch release automatically whenever a new commit lands
on `main`.

**Independent Test**: Push a new commit to `main` in a test repository where the
latest semver tag is `v0.4.2`, then verify the workflow computes `v0.4.3`, runs
validation first, and publishes tagged artifacts only after the build matrix
succeeds.

### Tests for User Story 1 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST and confirm they FAIL before implementation**

- [X] T008 [P] [US1] Add workflow contract tests for `push` triggering, `main` filtering, and publish-after-build ordering in `test/contract/release_workflow_contract_test.go`
- [X] T009 [P] [US1] Add integration tests for no-tag baselines and successive `main` pushes in `test/integration/release_workflow_test.go`
- [X] T010 [P] [US1] Add unit tests for patch-preview CLI defaults in `scripts/nextversion/main_test.go`

### Implementation for User Story 1

- [X] T011 [US1] Replace manual dispatch with automatic `push` triggering and `main`-only filtering in `.github/workflows/release.yml`
- [X] T012 [US1] Wire patch-only version calculation into `.github/workflows/release.yml` and `scripts/nextversion/main.go`
- [X] T013 [US1] Build and publish deterministic patch release artifacts in `.github/workflows/release.yml`
- [X] T014 [US1] Verify the Linux reference artifact directly with `ds --version` before publication in `.github/workflows/release.yml` and `test/contract/version_build_contract_test.go`
- [X] T015 [US1] Add workflow concurrency and `<15 minute` budget validation in `.github/workflows/release.yml` and `specs/002-versioning-release-automation/quickstart.md`

**Checkpoint**: User Story 1 publishes automatic patch releases from `main` and
is ready to validate as the MVP slice.

---

## Phase 4: User Story 2 - Keep automatic releases safe, idempotent, and actionable (Priority: P2)

**Goal**: Prevent reruns, failures, and non-semver repository noise from
creating duplicate or partial releases.

**Independent Test**: Rerun the workflow for an already released `main` commit,
then simulate failing tests or packaging and confirm no new tag or GitHub
release is created while logs explain the blocking condition.

### Tests for User Story 2 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST and confirm they FAIL before implementation**

- [X] T016 [P] [US2] Add workflow contract tests for already-released skips, secret-safe failure output, and the absence of manual release inputs in `test/contract/release_workflow_contract_test.go`
- [X] T017 [P] [US2] Add integration tests for already-tagged reruns, validation failures, and non-semver tag ignores in `test/integration/release_workflow_test.go`
- [X] T018 [P] [US2] Add unit tests for HEAD release detection and duplicate-tag guardrails in `internal/release/semver_test.go`

### Implementation for User Story 2

- [X] T019 [US2] Implement semver-tag-based already-released commit detection and duplicate-tag guardrails in `internal/release/semver.go`
- [X] T020 [US2] Add skip-on-rerun, fail-before-tag, and token-safe failure handling in `.github/workflows/release.yml`
- [X] T021 [US2] Align preview exit codes and workflow skip/failure vocabulary in `scripts/nextversion/main.go` and `.github/workflows/release.yml`
- [X] T022 [US2] Document rerun and failure behavior in `README.md` and `specs/002-versioning-release-automation/quickstart.md`

**Checkpoint**: User Story 2 makes the automatic release flow safe to rerun and
safe to fail.

---

## Phase 5: User Story 3 - Preview and verify automatic releases locally and in CI (Priority: P3)

**Goal**: Let contributors predict the next automatic patch release and verify
published binaries against the same metadata contract used in CI.

**Independent Test**: Run the local preview helper, build a binary with the
predicted release metadata, and compare its version output with the release
artifact contract used in CI.

### Tests for User Story 3 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST and confirm they FAIL before implementation**

- [X] T023 [P] [US3] Add contract tests for local next-patch preview and artifact self-reporting in `test/contract/version_build_contract_test.go`
- [X] T024 [P] [US3] Add integration tests for local preview parity and versioned artifact verification in `test/integration/version_build_test.go`
- [X] T025 [P] [US3] Add unit tests for version metadata fallback paths used by release verification in `pkg/dotenvsync/version_test.go`

### Implementation for User Story 3

- [X] T026 [US3] Document local preview, push-to-`main` release monitoring, and artifact verification in `README.md` and `specs/002-versioning-release-automation/quickstart.md`
- [X] T027 [US3] Align local preview examples with automatic patch semantics in `scripts/nextversion/main.go`, `README.md`, and `specs/002-versioning-release-automation/quickstart.md`
- [X] T028 [US3] Validate cross-platform release artifact version parity beyond the Linux reference check in `test/contract/version_build_contract_test.go`, `test/integration/version_build_test.go`, and `.github/workflows/release.yml`

**Checkpoint**: User Story 3 makes local and CI verification consistent and
independently usable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Finish documentation, regression coverage, UX consistency, and
budget validation across all stories.

- [X] T029 [P] Audit module boundaries and no-env-workflow impact in `.github/workflows/release.yml`, `internal/release/semver.go`, and `scripts/nextversion/main.go`
- [X] T030 [P] Audit consistent automatic-release vocabulary and secret-safe logs in `.github/workflows/release.yml`, `README.md`, and `test/contract/release_workflow_contract_test.go`
- [X] T031 Add regression coverage for concurrent `main` pushes and partial-release prevention in `test/contract/release_workflow_contract_test.go`, `test/integration/release_workflow_test.go`, and `internal/release/semver_test.go`
- [X] T032 Verify performance and dependency budgets in `.github/workflows/release.yml`, `scripts/nextversion/main.go`, and `go.mod`
- [X] T033 Run quickstart validation and refresh golden outputs in `specs/002-versioning-release-automation/quickstart.md`, `test/testdata/golden/release-published.txt`, and `test/testdata/golden/release-skipped.txt`

---

## Constitution Coverage

- **workflow-preservation**: T006-T007, T011-T015, T019-T022, T029-T030
- **module-boundaries**: T004-T007, T011-T013, T019-T021, T029, T032
- **deterministic-file-fidelity**: T001-T003, T013-T015, T023-T033
- **test-first-reliability**: T004-T005, T008-T010, T016-T018, T023-T025, T031
- **ux-consistency-secret-safety**: T008, T014, T016, T020-T022, T026-T030, T033
- **performance-dependency-budget**: T004-T007, T015, T024, T028, T032

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup and blocks all story work
- **User Story 1 (Phase 3)**: Depends on Foundational and is the recommended MVP
- **User Story 2 (Phase 4)**: Depends on User Story 1 because it hardens the same automatic release flow
- **User Story 3 (Phase 5)**: Depends on User Story 1 and reuses the published artifact flow for verification
- **Polish (Phase 6)**: Depends on all completed story phases

### User Story Dependencies

- **US1 (P1)**: First delivery target and recommended MVP
- **US2 (P2)**: Builds directly on the workflow introduced in US1
- **US3 (P3)**: Reuses US1 publication behavior and existing version metadata for verification

### Within Each User Story

- Tests must be written and fail before implementation
- Shared helpers precede workflow or preview wiring
- Workflow and preview implementation precede documentation and budget validation
- Each story must satisfy its independent test before moving to the next checkpoint

### Parallel Opportunities

- T002 and T003 can run in parallel after T001
- T004 and T005 can run in parallel before T006 and T007
- T008, T009, and T010 can run in parallel within US1
- T016, T017, and T018 can run in parallel within US2
- T023, T024, and T025 can run in parallel within US3
- T029 and T030 can run in parallel during polish

---

## Parallel Example: User Story 1

```bash
# Launch the failing tests for User Story 1 together:
Task: "T008 [US1] Add workflow contract tests in test/contract/release_workflow_contract_test.go"
Task: "T009 [US1] Add integration tests in test/integration/release_workflow_test.go"
Task: "T010 [US1] Add unit tests in scripts/nextversion/main_test.go"
```

## Parallel Example: User Story 2

```bash
# Launch the failing tests for User Story 2 together:
Task: "T016 [US2] Add workflow contract tests in test/contract/release_workflow_contract_test.go"
Task: "T017 [US2] Add integration tests in test/integration/release_workflow_test.go"
Task: "T018 [US2] Add semver unit tests in internal/release/semver_test.go"
```

## Parallel Example: User Story 3

```bash
# Launch the failing tests for User Story 3 together:
Task: "T023 [US3] Add contract tests in test/contract/version_build_contract_test.go"
Task: "T024 [US3] Add integration tests in test/integration/version_build_test.go"
Task: "T025 [US3] Add version metadata unit tests in pkg/dotenvsync/version_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **Stop and validate**: Confirm a push to `main` produces the next patch
   release after validation and build completion

### Incremental Delivery

1. Finish Setup and Foundational work
2. Deliver User Story 1 for automatic patch publication on `main`
3. Deliver User Story 2 for rerun safety and actionable failures
4. Deliver User Story 3 for local preview and CI verification parity
5. Finish polish, regression validation, and quickstart verification

### Parallel Team Strategy

1. One developer completes Setup and Foundational work
2. After Foundational:
   - Developer A implements User Story 1 in `.github/workflows/` and `scripts/nextversion/`
   - Developer B implements User Story 2 in `internal/release/` and workflow guardrails
   - Developer C implements User Story 3 in `README.md`, quickstart docs, and verification tests
3. Rejoin for polish, regression coverage, and budget validation

---

## Notes

- [P] tasks touch different files or become independent after shared prerequisites complete
- Each user story is independently testable at its checkpoint
- Suggested MVP scope: **User Story 1**
- The generated plan assumes release automation is patch-only on pushes to `main`; if major/minor automation is later required, refresh `spec.md`, `plan.md`, and this task list together
