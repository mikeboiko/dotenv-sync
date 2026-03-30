# dotenv-sync

`dotenv-sync` is a cross-platform Go CLI for keeping a local `.env` file aligned
with a schema in `.env.example` while resolving provider-managed values through
Bitwarden's `rbw` CLI.

The product name stays **dotenv-sync** and the default binary name is **`ds`**.

## Features

- `ds sync` writes `.env` from `.env.example` and `rbw`
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

```bash
VERSION=v0.1.0
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -o ./bin/ds \
  -ldflags "-X dotenv-sync/pkg/dotenvsync.Version=$VERSION -X dotenv-sync/pkg/dotenvsync.Commit=$COMMIT -X dotenv-sync/pkg/dotenvsync.BuildTime=$BUILD_TIME" \
  ./cmd/ds

./bin/ds --version
./bin/ds version
```

To preview the next release version from the current repository tags:

```bash
go run ./scripts/nextversion --bump patch
```

## CI

GitHub Actions runs `go test ./...` on every push via
`.github/workflows/go-tests.yml`.

GitHub Actions also supports a manual release workflow in
`.github/workflows/release.yml`. It must be triggered from the repository's
default branch head (typically `main`), computes the next semantic version from
existing `vX.Y.Z` tags, runs `go test ./...`, builds versioned archives, checks
the Linux release artifact with `ds --version`, and then publishes the GitHub
release.
