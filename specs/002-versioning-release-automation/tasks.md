---
description: 'Task list for version reporting and release automation'
---

# Tasks: Version reporting and release automation

**Input**: Design documents from `/specs/002-versioning-release-automation/`
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
- **Shared version metadata**: `pkg/dotenvsync/`
- **Release helpers**: `internal/release/`, `scripts/nextversion/`
- **GitHub automation**: `.github/workflows/`
- **Tests and fixtures**: `test/contract/`, `test/integration/`, `test/testdata/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish shared fixtures, helper extensions, and golden outputs for
version and release automation work.

- [ ] T001 Create versioning fixture and golden-output files in `test/testdata/golden/version-short.txt`, `test/testdata/golden/version-detailed.txt`, and `test/testdata/release/README.md`
- [ ] T002 [P] Extend shared metadata-aware command helpers in `test/contract/helpers_test.go` and `test/integration/helpers_test.go`
- [ ] T003 [P] Seed semver and artifact fixture cases in `test/testdata/release/mixed-tags.txt`, `test/testdata/release/no-tags.txt`, and `test/testdata/release/expected-assets.txt`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build shared metadata and release primitives required by the user
stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [ ] T004 [P] Add failing unit coverage for metadata defaults and render helpers in `pkg/dotenvsync/version_test.go`
- [ ] T005 [P] Add failing unit coverage for semver filtering, bumping, and artifact naming in `internal/release/semver_test.go`
- [ ] T006 Implement mutable build metadata variables and render helpers in `pkg/dotenvsync/version.go`
- [ ] T007 Implement semver filtering, next-version calculation, and artifact-name helpers in `internal/release/semver.go`

**Checkpoint**: Shared build metadata and release calculation primitives are
ready for story work.

---

## Phase 3: User Story 1 - Inspect the installed CLI version (Priority: P1) 🎯 MVP

**Goal**: Let users inspect the installed `ds` build quickly with `--version`
and in detail with `ds version`.

**Independent Test**: Build `ds` once as a development binary and once with
explicit `ldflags`, then verify `ds --version` and `ds version` print the
expected metadata without touching `.env` files or provider state.

### Tests for User Story 1 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T008 [P] [US1] Add CLI contract coverage for `ds --version` and `ds version` in `test/contract/version_contract_test.go`
- [ ] T009 [P] [US1] Add integration coverage for development and ldflags-injected builds in `test/integration/version_command_test.go`
- [ ] T010 [P] [US1] Add root and version command unit coverage in `internal/cli/root_test.go` and `internal/cli/version_test.go`

### Implementation for User Story 1

- [ ] T011 [US1] Implement concise and detailed version renderers in `internal/cli/version.go`
- [ ] T012 [US1] Wire `--version` handling and register the `version` subcommand in `internal/cli/root.go`
- [ ] T013 [US1] Align `ds version` output with shared metadata helpers in `internal/cli/version.go` and `pkg/dotenvsync/version.go`
- [ ] T014 [US1] Enforce no-extra-args and success exit behavior for version paths in `internal/cli/version.go` and `internal/cli/root.go`
- [ ] T015 [US1] Verify fast-path version performance in `internal/cli/version_benchmark_test.go` and `test/contract/version_contract_test.go`

**Checkpoint**: User Story 1 is independently functional and ready to demo as
the MVP slice.

---

## Phase 4: User Story 2 - Publish semver releases from GitHub Actions (Priority: P2)

**Goal**: Let maintainers publish semver releases through a manual GitHub
Actions workflow that computes the next version, validates the repo, and
publishes tagged artifacts.

**Independent Test**: Exercise semver bump inputs against fixture repos and
workflow expectations, then verify that release publication occurs only after
tests and builds succeed.

### Tests for User Story 2 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T016 [P] [US2] Add contract coverage for workflow dispatch inputs, publish ordering, and artifact naming in `test/contract/release_workflow_contract_test.go`
- [ ] T017 [P] [US2] Add integration coverage for semver bump selection, no-tag baselines, and duplicate-tag rejection in `test/integration/release_workflow_test.go`
- [ ] T018 [P] [US2] Add unit coverage for latest-tag lookup and default-branch enforcement in `internal/release/semver_test.go`

### Implementation for User Story 2

- [ ] T019 [US2] Extend `internal/release/semver.go` with latest-tag lookup, duplicate-tag rejection, and default-branch validation
- [ ] T020 [US2] Implement the release preview CLI in `scripts/nextversion/main.go`
- [ ] T021 [US2] Create the manual semver release workflow in `.github/workflows/release.yml`
- [ ] T022 [US2] Validate publish-after-build ordering and actionable workflow failure messages in `.github/workflows/release.yml` and `test/contract/release_workflow_contract_test.go`

**Checkpoint**: User Story 2 can publish tagged releases with repeatable semver
calculation and guarded workflow behavior.

---

## Phase 5: User Story 3 - Build repeatable versioned binaries locally and in CI (Priority: P3)

**Goal**: Let contributors build and verify versioned binaries locally and in CI
using the same metadata contract as published releases.

**Independent Test**: Follow the documented local `ldflags` flow, run `ds
version`, and confirm the output matches the release artifact contract and CI
verification steps.

### Tests for User Story 3 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T023 [P] [US3] Add contract coverage for documented local ldflags builds and artifact self-reporting in `test/contract/version_build_contract_test.go`
- [ ] T024 [P] [US3] Add integration coverage for local ldflags builds and release-artifact version consistency in `test/integration/version_build_test.go`
- [ ] T025 [P] [US3] Add unit coverage for platform reporting and fallback metadata in `pkg/dotenvsync/version_test.go`

### Implementation for User Story 3

- [ ] T026 [US3] Add release-artifact self-verification steps to `.github/workflows/release.yml` using `ds --version`
- [ ] T027 [US3] Document local ldflags builds, version inspection, and maintainer release steps in `README.md` and `specs/002-versioning-release-automation/quickstart.md`
- [ ] T028 [US3] Align local release-preview examples with `scripts/nextversion/main.go` and `specs/002-versioning-release-automation/quickstart.md`

**Checkpoint**: User Story 3 makes local, CI, and release artifact verification
consistent and independently usable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Finish documentation, regression coverage, UX consistency, and
budget validation across all stories.

- [ ] T029 [P] Update help text and installation documentation for version inspection in `README.md` and `internal/cli/version.go`
- [ ] T030 [P] Audit consistent version vocabulary and token-safe release logs in `internal/cli/version.go`, `.github/workflows/release.yml`, and `specs/002-versioning-release-automation/contracts/cli-version.contract.md`
- [ ] T031 Add regression coverage for non-semver tags and partial-release prevention in `test/contract/release_workflow_contract_test.go`, `test/integration/release_workflow_test.go`, and `internal/release/semver_test.go`
- [ ] T032 Verify performance and dependency budgets in `internal/cli/version_benchmark_test.go` and `.github/workflows/release.yml`
- [ ] T033 Run quickstart validation and refresh golden outputs in `specs/002-versioning-release-automation/quickstart.md`, `test/testdata/golden/version-short.txt`, and `test/testdata/golden/version-detailed.txt`

---

## Constitution Coverage

- **workflow-preservation**: T008-T014, T016-T022
- **module-boundaries**: T004-T007, T011-T013, T019-T021
- **deterministic-file-fidelity**: T001-T003, T005, T015-T018, T021-T028, T033
- **test-first-reliability**: T004-T005, T008-T010, T016-T018, T023-T025, T031
- **ux-consistency-secret-safety**: T008-T015, T016, T022, T027-T030
- **performance-dependency-budget**: T005-T007, T015, T021, T026, T032

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion - recommended MVP slice
- **User Story 2 (Phase 4)**: Depends on Foundational completion and reuses shared metadata/release helpers, but remains independently testable from User Story 1
- **User Story 3 (Phase 5)**: Depends on Foundational completion and validates the same release metadata contract locally and in CI
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: First delivery target and recommended MVP
- **User Story 2 (P2)**: Can proceed after Foundational, but is easiest to validate once User Story 1 exposes visible version metadata
- **User Story 3 (P3)**: Can proceed after Foundational and benefits from the workflow created in User Story 2

### Within Each User Story

- Tests MUST be written and fail before implementation
- Shared helpers precede command or workflow wiring
- CLI or workflow wiring precedes documentation and budget validation
- Each story must satisfy its own independent test before moving on

### Parallel Opportunities

- T002-T003 can run in parallel after T001
- T004-T005 can run in parallel before T006-T007
- Within **US1**, T008-T010 can run in parallel before T011-T015
- Within **US2**, T016-T018 can run in parallel before T019-T022
- Within **US3**, T023-T025 can run in parallel before T026-T028
- T029-T030 can run in parallel during polish once the stories are complete

---

## Parallel Example: User Story 1

```bash
# Launch the required tests for User Story 1 together:
Task: "T008 [US1] Add CLI contract coverage in test/contract/version_contract_test.go"
Task: "T009 [US1] Add integration coverage in test/integration/version_command_test.go"
Task: "T010 [US1] Add unit coverage in internal/cli/root_test.go and internal/cli/version_test.go"
```

## Parallel Example: User Story 2

```bash
# Launch the required tests for User Story 2 together:
Task: "T016 [US2] Add workflow contract coverage in test/contract/release_workflow_contract_test.go"
Task: "T017 [US2] Add release integration coverage in test/integration/release_workflow_test.go"
Task: "T018 [US2] Add semver unit coverage in internal/release/semver_test.go"
```

## Parallel Example: User Story 3

```bash
# Launch the required tests for User Story 3 together:
Task: "T023 [US3] Add local-build contract coverage in test/contract/version_build_contract_test.go"
Task: "T024 [US3] Add local-build integration coverage in test/integration/version_build_test.go"
Task: "T025 [US3] Add metadata fallback unit coverage in pkg/dotenvsync/version_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Confirm `ds --version` and `ds version` work for both
   dev and ldflags-injected builds

### Incremental Delivery

1. Finish Setup + Foundational
2. Deliver User Story 1 for immediate operator value
3. Add User Story 2 for maintainer release automation
4. Add User Story 3 for local and CI verification parity
5. Finish polish and regression validation

### Parallel Team Strategy

1. One developer completes Setup + Foundational
2. After Foundational:
   - Developer A owns User Story 1 CLI version surfaces
   - Developer B owns User Story 2 release automation
   - Developer C owns User Story 3 local/CI verification and documentation
3. Rejoin for polish, regression coverage, and quickstart validation

---

## Notes

- [P] tasks use different files or can proceed after shared prerequisites complete
- User story phases map directly to independently testable feature slices
- Suggested MVP scope: **User Story 1 only**
