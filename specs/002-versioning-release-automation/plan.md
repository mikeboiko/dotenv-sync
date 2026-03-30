# Implementation Plan: Version reporting and release automation

**Branch**: `002-versioning-release-automation` | **Date**: 2026-03-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-versioning-release-automation/spec.md`

## Summary

Add trustworthy version reporting to `ds` through a fast `--version` flag and a
detailed `version` command, then add a maintainer-triggered GitHub Actions
workflow that bumps semantic versions, tags releases, and publishes versioned
cross-platform binaries without manual source edits. The design keeps existing
dotenv workflows untouched, uses build-time metadata injection for the binary,
and keeps release logic testable and isolated from sync/provider code.

## Technical Context

**Language/Version**: Go 1.22  
**Primary Dependencies**: Existing `github.com/spf13/cobra`, Go standard
library, GitHub Actions `actions/checkout` and `actions/setup-go`, GitHub CLI
`gh` on release runners  
**Storage**: Local source files, Git tags, and GitHub Releases metadata  
**Testing**: `go test ./...`, CLI contract tests, subprocess integration tests,
unit tests for semver/version formatting, and workflow smoke validation  
**Target Platform**: Cross-platform CLI for Linux, macOS, and Windows; release
automation on GitHub-hosted Linux runners  
**Project Type**: Single-binary CLI with GitHub-hosted release automation  
**Performance Goals**: `ds --version` and `ds version` complete within 100 ms
locally; release workflow completes validation, cross-build, and publication for
the supported matrix within 15 minutes under normal runner availability  
**Constraints**: No manual version edits in Go source per release, no runtime
network calls for version inspection, no new runtime dependencies in `ds`,
existing `go test ./...` remains the release gate, release truth comes from
semantic Git tags, and workflow logs must remain token-safe  
**Scale/Scope**: One repository, one public CLI binary (`ds`), one manual
release workflow, and a supported release matrix of Linux amd64/arm64, macOS
amd64/arm64, and Windows amd64 artifacts

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

Pre-design review passes because the feature adds version visibility and release
automation without changing the existing `.env` workflow or broadening runtime
scope. Post-design review also passes after `research.md`, `data-model.md`,
`quickstart.md`, and the contracts below define testable CLI behavior, release
guardrails, and minimal-dependency publication paths.

| Check ID                      | Status | Evidence                                                                                                                                                                                                                         |
| ----------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| workflow-preservation         | PASS   | The Summary and Technical Context keep release/version behavior separate from sync flows, and `contracts/cli-version.contract.md` adds inspection-only commands with no `.env` mutation or runtime wrappers.                     |
| module-boundaries             | PASS   | Project Structure isolates CLI wiring in `internal/cli`, embedded metadata in `pkg/dotenvsync`, semver release logic in `internal/release`, and workflow-only glue in `scripts/nextversion` and `.github/workflows/release.yml`. |
| deterministic-file-fidelity   | PASS   | The feature does not alter envfile mutation semantics, and `contracts/release-workflow.contract.md` plus `quickstart.md` require deterministic artifact naming and no-op-safe version inspection paths.                          |
| test-first-reliability        | PASS   | Technical Context requires unit, contract, integration, and workflow validation coverage before implementation, including semver bump edge cases and development-build fallbacks.                                                |
| ux-consistency-secret-safety  | PASS   | `spec.md` UX requirements define stable version output, actionable workflow failures, and token-safe logs, while the CLI contract standardizes concise and detailed version output.                                              |
| performance-dependency-budget | PASS   | Technical Context keeps version inspection constant-time, avoids new runtime dependencies, and uses existing GitHub tooling (`go build`, `gh`) rather than adding heavier release frameworks.                                    |

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
├── cli/
│   ├── root.go
│   └── version.go
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
│   └── version_contract_test.go
└── integration/
    ├── version_command_test.go
    └── release_workflow_test.go

README.md
```

**Structure Decision**: Keep user-facing version reporting in the existing CLI
package, centralize embedded build metadata in `pkg/dotenvsync`, place
maintainer-only semver calculation in a focused `internal/release` package, and
use a small repository-local helper under `scripts/nextversion` so GitHub
Actions can compute release versions without leaking release plumbing into
provider or envfile packages.

## Complexity Tracking

No constitution violations or complexity exceptions are expected for this
feature.
