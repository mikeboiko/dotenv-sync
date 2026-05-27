---
description: Project-specific Copilot guidance derived from active feature plans.
---

# dotenv-sync Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-05-27

## Active Technologies

- Go 1.22 + `github.com/spf13/cobra` for CLI routing, `gopkg.in/yaml.v3` for config loading, Go standard-library JSON support, and provider CLIs such as `rbw` and `keepassxc-cli` (003-ds-push, 001-dotenv-sync-cli)
- Local env files plus provider-managed secrets, including a repo-scoped Bitwarden login-item notes payload in `note_json` mode for `ds push` (003-ds-push)

- Go 1.22
- `github.com/spf13/cobra` for CLI routing (001-dotenv-sync-cli)
- `gopkg.in/yaml.v3` for optional `.envsync.yaml` configuration (001-dotenv-sync-cli)
- Provider adapters under `internal/provider/`, currently including Bitwarden via `rbw` and KeePass via `keepassxc-cli` (001-dotenv-sync-cli)
- Git tags and GitHub Releases as the versioning source of truth (002-versioning-release-automation)
- GitHub Actions automatic patch-release workflow for `ds` with downstream package publishers such as AUR (002-versioning-release-automation)

## Project Structure

```text
cmd/ds/
internal/cli/
internal/config/
internal/envfile/
internal/fs/
internal/provider/bitwarden/
internal/provider/keepass/
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

- 003-ds-push: Added Bitwarden write-back via `ds push`; push remains Bitwarden-only while `sync`, `diff`, `validate`, `doctor`, and related flows support multiple providers

- 002-versioning-release-automation: Planned `ds --version`, `ds version`, build-time metadata injection, and manual GitHub Actions semver releases

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
