# Implementation Plan: Automatic patch release automation

**Branch**: `002-versioning-release-automation` | **Date**: 2026-03-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-versioning-release-automation/spec.md`

## Summary

Convert the existing manual release flow into an automatic patch-release pipeline
that runs on every push to `main`. The design keeps the current dotenv workflow
untouched, reuses the existing embedded version metadata for artifact
verification, centralizes patch calculation and rerun safety in
`internal/release`, and keeps workflow-only orchestration in
`.github/workflows/release.yml`.

## Technical Context

**Language/Version**: Go 1.22
**Primary Dependencies**: Go standard library, existing `github.com/spf13/cobra`
for the already-shipped version surfaces, GitHub Actions `actions/checkout` and
`actions/setup-go`, GitHub CLI `gh` on release runners
**Storage**: Local source files, Git tags, GitHub Releases metadata, and GitHub
workflow event data
**Testing**: `go test ./...`, workflow contract tests, integration tests for
patch preview and rerun guards, and unit tests for semver/release-state helpers
**Target Platform**: Cross-platform CLI for Linux, macOS, and Windows; automatic
release automation on GitHub-hosted Linux runners
**Project Type**: Single-binary CLI with push-triggered GitHub release
automation
**Performance Goals**: Patch preview completes within 1 second locally; the
release workflow completes validation, cross-build, and publication for the
supported matrix within 15 minutes under normal runner availability
**Constraints**: No manual bump selection, no manual source edits or local tag
creation, existing `go test ./...` remains the release gate, release truth comes
from semantic Git tags, only one release may be published per `main` commit, and
workflow logs must remain token-safe
**Scale/Scope**: One repository, one default branch (`main`), patch-only release
automation, and a supported release matrix of Linux amd64/arm64, macOS
amd64/arm64, and Windows amd64 artifacts

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

Pre-design review passes because the feature changes only maintainer-facing
release automation and leaves `.env`, `.env.example`, and provider flows
untouched. Post-design review also passes after `research.md`, `data-model.md`,
`quickstart.md`, and the release contract define testable patch-only automation,
race-safe publication, and secret-safe operator messaging.

| Check ID                      | Status | Evidence |
| ----------------------------- | ------ | -------- |
| workflow-preservation         | PASS   | The Summary and Technical Context keep all changes in release automation files and reuse existing version metadata only for verification, with no `.env` or provider-path changes. |
| module-boundaries             | PASS   | Project Structure isolates semver and release-state logic in `internal/release`, local preview glue in `scripts/nextversion`, and workflow orchestration in `.github/workflows/release.yml`. |
| deterministic-file-fidelity   | PASS   | The feature does not alter envfile mutation semantics, and the workflow contract requires deterministic artifact names, publish-after-build ordering, and skip-without-mutation reruns. |
| test-first-reliability        | PASS   | Technical Context requires failing workflow contract, integration, and unit coverage before implementation, including rerun, duplicate-tag, and failure-path cases. |
| ux-consistency-secret-safety  | PASS   | `spec.md` and `contracts/release-workflow.contract.md` define consistent skip/publish/failure messaging and require logs to expose only refs, versions, and artifact names. |
| performance-dependency-budget | PASS   | Patch-only automation reuses existing GitHub tooling and version metadata, avoids new runtime dependencies, and keeps the workflow within the existing release budget. |

## Project Structure

### Documentation (this feature)

```text
specs/002-versioning-release-automation/
├── plan.md
├── spec.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cli-version.contract.md
│   └── release-workflow.contract.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/
└── workflows/
    ├── go-tests.yml
    └── release.yml

cmd/
└── ds/
    └── main.go

internal/
└── release/
    └── semver.go

pkg/
└── dotenvsync/
    └── version.go

scripts/
└── nextversion/
    └── main.go

test/
├── contract/
│   ├── release_workflow_contract_test.go
│   └── version_build_contract_test.go
└── integration/
    ├── release_workflow_test.go
    └── version_build_test.go

README.md
```

**Structure Decision**: Keep patch calculation, tag filtering, branch validation,
and rerun safety in `internal/release`; keep local preview and operator-facing
CLI glue in `scripts/nextversion`; keep workflow-specific checkout, concurrency,
build, and publication steps in `.github/workflows/release.yml`; and reuse the
existing embedded metadata in `pkg/dotenvsync` only for artifact verification.

## Complexity Tracking

No constitution violations or complexity exceptions are expected for this
feature.
