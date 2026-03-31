# dotenv-sync

`dotenv-sync` is a cross-platform Go CLI for keeping a local `.env` file aligned
with a schema in `.env.example` while resolving provider-managed values through
Bitwarden's `rbw` CLI.

The product name stays **dotenv-sync** and the default binary name is **`ds`**.

## Features

- `ds sync` writes `.env` from `.env.example` and `rbw`
- `ds push` uploads the current `.env` into a repo-scoped Bitwarden item
- `ds diff` previews drift without writing files
- `ds validate` reports malformed files, drift, duplicates, and missing secrets
- `ds doctor` checks config and `rbw` readiness
- `ds init` bootstraps `.env.example` from `.env`
- `ds missing` lists unresolved provider-backed keys
- `ds reverse` adds missing schema placeholders back into `.env.example`
- `ds --version` and `ds version` report build and release metadata

## Install and build

```bash
go build -o ./bin/ds ./cmd/ds
```

Local development builds report `dev` metadata by default. Release builds inject
their version, commit, and build time at build time instead of editing source
files.

## Keep `ds` up to date locally

If you keep `ds` in `~/.local/bin/ds`, the simplest update flow is:

```bash
git switch main
git pull --ff-only
./scripts/install-local.sh
```

By default `./scripts/install-local.sh` installs to `~/.local/bin/ds`. You can
override that with `--bin /custom/path/to/ds`.

If you use this repository's `lefthook` setup, `lefthook install` also enables
automatic local refreshes after `git commit` and `git merge` on the default
branch. Those hooks call `./scripts/install-local.sh --quiet` through
`./scripts/install-local-hook.sh`, and they intentionally skip non-default
branches so feature work does not overwrite your globally installed `ds`.

The script injects version metadata from your current checkout:

- exact release tag checkout: `v0.4.0`
- commits ahead of a release: `v0.4.0-3-gabc1234`
- no release tags yet: `dev-abc1234`

That means:

- if `main` is exactly on the latest release tag, `ds version` matches what a
  GitHub release installer sees
- if `main` is ahead of the latest release, your local build intentionally shows
  a newer git-derived version so it does not pretend to be the last published
  release

If you want the exact same version string as the latest GitHub release, build
from the release tag or install the GitHub release artifact.

## Add `ds` to your `PATH`

You can always run the binary directly as `./bin/ds` (or `bin\ds.exe` on
Windows), but for a normal `ds ...` workflow you can either add the build
directory to your `PATH` or copy/symlink the binary into a directory that is
already on your `PATH`.

### POSIX shells (`bash`, `zsh`, `sh`)

Add your build directory to a shell startup file such as `~/.profile`,
`~/.bashrc`, or `~/.zshrc`:

```bash
export PATH="$PATH:/absolute/path/to/dotenv-sync/bin"
```

Reload the file or open a new shell:

```bash
source ~/.bashrc
```

### `fish`

```fish
fish_add_path /absolute/path/to/dotenv-sync/bin
```

### PowerShell

For the current session:

```powershell
$env:Path += ";C:\absolute\path\to\dotenv-sync\bin"
```

For a persistent install, add the same directory through your system PATH
settings or place `ds.exe` in a directory that is already on PATH.

### Alternative: copy or symlink the binary

Examples:

```bash
ln -s /absolute/path/to/dotenv-sync/bin/ds ~/.local/bin/ds
```

or:

```bash
install -Dm755 ./bin/ds ~/.local/bin/ds
```

## Configuration

`.envsync.yaml` is optional:

```yaml
provider: bitwarden
schema_file: .env.example
env_file: .env
item_name: my-app
storage_mode: fields
mapping:
  DATABASE_URL: db_url
  JWT_SECRET: auth_jwt
```

Blank values in `.env.example` are treated as provider-managed secrets. Literal
values are treated as safe defaults and copied into `.env`.

If `item_name` is omitted, `ds` derives it from the Git repository root
directory name and falls back to the current working directory name when Git
metadata is unavailable. By default, provider-managed keys resolve as
`rbw get <item_name> --field <ENV_VAR>`, and `mapping` overrides only the field
name inside that Bitwarden item.

`storage_mode` defaults to `fields` for backward-compatible reads from the
repo-scoped Bitwarden item fields. In `fields` mode, `ds push` can update
provider-managed keys that map to Bitwarden's built-in `password` field. Set
`storage_mode: note_json` to store the full repo env map in the item notes for
round-trip `push`/`sync` workflows.

That makes shared aliases across repos possible with the default field-based
layout. For example, both repos below read and write the same Bitwarden value:

```yaml
# repo1/.envsync.yaml
item_name: shared-dev
mapping:
  DB_PASSWD: password
```

```yaml
# repo2/.envsync.yaml
item_name: shared-dev
mapping:
  PSWD: password
```

## Commands

### `ds sync`

```bash
ds sync
ds sync --dry-run
```

- Reads `.env.example` as the schema contract
- Resolves blank entries through `rbw` using a repo-scoped item by default
- Preserves comment order and line endings when rewriting `.env`
- Produces `WRITTEN`, `UNCHANGED`, and `MISSING` output vocabulary for sync runs
- Uses `CHECKED` summaries for dry-run previews

### `ds push`

```bash
ds push --dry-run
ds push
```

- Requires `storage_mode: note_json`
- Reads `.env` as the upload source and `.env.example` as schema context
- In `note_json`, writes a deterministic JSON payload into the repo-scoped
  Bitwarden item notes
- In `fields`, updates provider-managed keys present in `.env` when they map to
  Bitwarden's built-in `password` field
- Never prints raw values; previews use redacted markers such as `[REDACTED]`
- Leaves `.env` and `.env.example` untouched

Field mode is useful for shared aliases across repos, but `rbw` cannot script
writes to arbitrary custom Bitwarden fields. If a repo maps pushed keys to
custom fields, `ds push` fails with actionable guidance to use `password` or
switch to `storage_mode: note_json`.

### `ds diff`

```bash
ds diff
```

Prints only real changes using redacted markers such as `[RESOLVED]` and
`[STATIC]`.

### `ds validate`

```bash
ds validate
```

Returns exit code `2` when drift, malformed input, duplicates, or unresolved
secrets are found.

### `ds doctor`

```bash
ds doctor
```

Checks `.envsync.yaml` readability and `rbw` readiness without printing any
secret values.

### `ds init`

```bash
ds init
ds init --dry-run
```

Creates `.env.example` from `.env`, blanking secret-like values while copying
safe defaults.

### `ds missing`

```bash
ds missing
```

Lists only unresolved provider-managed schema keys and exits with code `2` when
any key is missing.

### `ds reverse`

```bash
ds reverse --dry-run
ds reverse
```

Adds keys found in `.env` but missing from `.env.example` as blank placeholders.

### `ds --version` and `ds version`

```bash
ds --version
ds version
```

- `ds --version` prints a concise version string such as `ds v0.4.0`
- `ds version` prints detailed metadata including version, commit, build time,
  and platform
- Local development builds fall back to `dev`, `none`, and `unknown`

## Exit codes

- `0`: success or no-op success
- `1`: operational failure
- `2`: validation, drift, or missing-value issue

## Development

```bash
go test ./...
go test ./... -run TestContract
go test ./... -bench . -run '^$'
```

## Build a local versioned binary

Preview the next automatic patch release from the current reachable semver tags:

```bash
go run ./scripts/nextversion
```

Example outputs:

```text
v0.0.1
v0.4.3
```

Then build a local binary with that predicted release metadata:

```bash
VERSION=$(go run ./scripts/nextversion)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -o ./bin/ds \
  -ldflags "-X dotenv-sync/pkg/dotenvsync.Version=$VERSION -X dotenv-sync/pkg/dotenvsync.Commit=$COMMIT -X dotenv-sync/pkg/dotenvsync.BuildTime=$BUILD_TIME" \
  ./cmd/ds

./bin/ds --version
./bin/ds version
```

Local development builds without injected metadata fall back to `dev`, `none`,
and `unknown`.

To install straight into `~/.local/bin/ds` with Git-derived version metadata:

```bash
./scripts/install-local.sh
```

## CI

GitHub Actions runs `go test ./...` on every push via
`.github/workflows/go-tests.yml`.

GitHub Actions also runs `.github/workflows/release.yml` automatically on every
push to `main`. The workflow computes the next patch version from reachable
`vX.Y.Z` tags, skips reruns when the pushed commit is already released by a
semver tag, runs `go test ./...`, builds versioned archives, verifies the Linux
reference artifact with `ds --version`, and then publishes the GitHub release.

You can monitor the latest automatic release run with:

```bash
gh run list --workflow release.yml --limit 1
gh run watch <run-id>
```

If a rerun finds that the commit already has a release tag, the workflow exits
without publishing again. If the tag exists but the GitHub release record is
missing, repair that release manually instead of expecting the workflow to
recreate it.
