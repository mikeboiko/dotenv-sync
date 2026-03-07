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

## Install and build

```bash
go build -o ./bin/ds ./cmd/ds
```

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
mapping:
  DATABASE_URL: database-url
  JWT_SECRET: jwt-secret
```

Blank values in `.env.example` are treated as provider-managed secrets. Literal
values are treated as safe defaults and copied into `.env`.

## Commands

### `ds sync`

```bash
ds sync
ds sync --dry-run
```

- Reads `.env.example` as the schema contract
- Resolves blank entries through `rbw`
- Preserves comment order and line endings when rewriting `.env`
- Produces `CHECKED`, `UNCHANGED`, `WRITTEN`, and `MISSING` output vocabulary

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
