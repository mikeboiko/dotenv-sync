---
description: Project-specific Copilot guidance derived from active feature plans.
---

# dotenv-sync Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-07

## Active Technologies

- Go 1.22
- `github.com/spf13/cobra` for CLI routing (001-dotenv-sync-cli)
- `gopkg.in/yaml.v3` for optional `.envsync.yaml` configuration (001-dotenv-sync-cli)
- `rbw` CLI as the Bitwarden runtime prerequisite (001-dotenv-sync-cli)

## Project Structure

```text
cmd/ds/
internal/cli/
internal/config/
internal/envfile/
internal/fs/
internal/provider/bitwarden/
internal/report/
internal/sync/
pkg/dotenvsync/
test/contract/
test/integration/
test/testdata/
```

## Commands

- `go test ./...`
- `go test ./... -run TestContract`
- `go test ./... -bench .`

## Code Style

Go 1.22: Follow standard conventions

## Recent Changes

- 001-dotenv-sync-cli: Planned a cross-platform Go CLI with deterministic envfile handling, `rbw`-backed Bitwarden integration, YAML config, the `ds` binary, and CLI contract coverage

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
