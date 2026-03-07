# Implementation Plan: dotenv-sync CLI MVP

**Branch**: `001-dotenv-sync-cli` | **Date**: 2026-03-07 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-dotenv-sync-cli/spec.md`

## Summary

Build `dotenv-sync` as a cross-platform Go CLI that keeps `.env.example` as the
schema contract and produces trustworthy local `.env` files without changing the
standard developer workflow. The design uses Cobra-based commands, a
fidelity-preserving envfile parser and writer, a provider interface with a
Bitwarden CLI adapter, and shared reporting and redaction utilities so sync,
diff, validate, doctor, init, missing, and reverse-sync behaviors stay
consistent across Linux, macOS, and Windows.

## Technical Context

**Language/Version**: Go 1.26  
**Primary Dependencies**: `github.com/spf13/cobra` for CLI routing,
`gopkg.in/yaml.v3` for optional `.envsync.yaml` config, Go standard library for
file I/O, JSON parsing, process execution, and testing  
**Storage**: Local files only (`.env.example`, `.env`, `.envsync.yaml`)  
**Testing**: `go test`, table-driven unit tests, golden-file tests,
subprocess-based CLI integration tests, benchmark tests for sync, diff, and
validate hot paths  
**Target Platform**: Cross-platform CLI for Linux, macOS, and Windows  
**Project Type**: Single-binary CLI application  
**Performance Goals**: <=200 ms p95 for local parse, diff, and validate on
500-key files, <=1 s no-op sync detection, <=10 s provider-backed sync and
validate with at most one provider lookup per distinct key per command  
**Constraints**: No runtime command wrapping, deterministic file rewrites with
comment, order, and line-ending preservation, secret-safe output, minimal
runtime dependencies, Bitwarden as the first provider, cross-platform path and
process handling, CI-friendly exit codes  
**Scale/Scope**: One project directory per invocation, up to 500 keys per
schema, one provider shipping initially, local developer and CI workflows only
for MVP

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Pre-design review passes from the accepted spec and user requirement for a
cross-platform Go implementation. Post-design re-check also passes after the
research, data model, quickstart, and contract artifacts below were produced.

| Check ID | Status | Evidence |
|----------|--------|----------|
| workflow-preservation | PASS | Summary and Technical Context keep `.env.example` as the schema and preserve direct `.env` usage with no runtime wrapper commands. `research.md` Decisions 1, 4, and 5 keep the workflow file-centric, while `contracts/cli-commands.contract.md` defines first-class commands instead of wrapped execution. |
| module-boundaries | PASS | Project Structure isolates command routing in `internal/cli`, env parsing and writing in `internal/envfile`, sync orchestration in `internal/sync`, provider adapters in `internal/provider`, configuration in `internal/config`, and output and redaction in `internal/report`. `data-model.md` separates schema, provider, and sync-result entities so logic does not bleed across packages. |
| deterministic-file-fidelity | PASS | Technical Context requires deterministic rewrites and preserved line endings. `research.md` Decision 3 commits to a custom token-based parser and writer, and `contracts/env-file-format.contract.md` defines the round-trip and no-op guarantees the implementation must preserve. |
| test-first-reliability | PASS | Technical Context requires unit, integration, golden, and benchmark tests. `research.md` Decision 7 defines failing-first automated coverage, `quickstart.md` includes the validation commands to run, and the planned structure includes `test/contract/`, `test/integration/`, and package-level `_test.go` files for negative-path coverage. |
| ux-consistency-secret-safety | PASS | The spec's UX requirements are carried into the design through shared reporting and redaction components in Project Structure, `contracts/output-format.contract.md`, and `contracts/error-codes.contract.md`. The plan explicitly standardizes status vocabulary, recovery guidance, and redacted output across all commands. |
| performance-dependency-budget | PASS | Technical Context sets the latency and lookup budgets, while `research.md` Decisions 2, 3, 4, and 6 justify the minimal dependency set, one-lookup-per-key strategy, and cross-platform file-write behavior. The structure keeps performance-sensitive logic in focused packages that can be benchmarked directly. |

## Project Structure

### Documentation (this feature)

```text
specs/001-dotenv-sync-cli/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cli-commands.contract.md
│   ├── env-file-format.contract.md
│   ├── output-format.contract.md
│   ├── error-codes.contract.md
│   └── provider-adapter.contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── dotenv-sync/
    └── main.go

internal/
├── cli/
│   ├── root.go
│   ├── sync.go
│   ├── diff.go
│   ├── validate.go
│   ├── doctor.go
│   ├── init.go
│   ├── missing.go
│   └── reverse.go
├── config/
│   └── config.go
├── envfile/
│   ├── model.go
│   ├── parser.go
│   ├── writer.go
│   └── merge.go
├── fs/
│   └── atomic.go
├── provider/
│   ├── provider.go
│   └── bitwarden/
│       ├── adapter.go
│       ├── client.go
│       └── status.go
├── report/
│   ├── output.go
│   ├── redact.go
│   └── exitcodes.go
└── sync/
    ├── engine.go
    ├── forward.go
    ├── reverse.go
    ├── diff.go
    ├── validate.go
    └── missing.go

pkg/
└── dotenvsync/
    └── version.go

test/
├── contract/
├── integration/
└── testdata/
    ├── env/
    ├── provider/
    └── golden/
```

**Structure Decision**: Use a single Go module with one CLI entry point and
internal packages for envfile fidelity, provider adapters, orchestration,
reporting, and configuration. This keeps the cross-platform binary cohesive,
aligns with the user's recommended Go layout, and avoids unnecessary module or
service boundaries for an MVP focused on local developer workflows.

## Complexity Tracking

No constitution violations or complexity exceptions require justification at
this stage.
