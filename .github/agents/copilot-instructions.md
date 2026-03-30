---
description: Project-specific Copilot guidance derived from active feature plans.
---

# dotenv-sync Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-30

## Active Technologies

- Go 1.22 + `github.com/spf13/cobra` for CLI routing, `gopkg.in/yaml.v3` for config loading, Go standard-library JSON support, installed `rbw` CLI (003-ds-push)
- Local env files plus a repo-scoped Bitwarden login-item notes payload in `note_json` mode (003-ds-push)

- Go 1.22
- `github.com/spf13/cobra` for CLI routing (001-dotenv-sync-cli)
- `gopkg.in/yaml.v3` for optional `.envsync.yaml` configuration (001-dotenv-sync-cli)
- `rbw` CLI as the Bitwarden runtime prerequisite (001-dotenv-sync-cli)
- Git tags and GitHub Releases as the versioning source of truth (002-versioning-release-automation)
- GitHub Actions manual semver release workflow for `ds` (002-versioning-release-automation)

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

- 003-ds-push: Added Go 1.22 + `github.com/spf13/cobra` for CLI routing, `gopkg.in/yaml.v3` for config loading, Go standard-library JSON support, installed `rbw` CLI

- 002-versioning-release-automation: Planned `ds --version`, `ds version`, build-time metadata injection, and manual GitHub Actions semver releases

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
