# Quickstart: dotenv-sync CLI MVP

## Goal

Validate the planned cross-platform Go CLI workflow for syncing `.env` files
from `.env.example` and Bitwarden without changing normal developer commands.

## Prerequisites

1. Go 1.26 or newer is installed.
2. Bitwarden CLI (`bw`) is installed and available on `PATH`.
3. The user is logged in to Bitwarden and can unlock the target vault.
4. The repository contains the planned source tree and command implementations.

## 1. Build the CLI

```bash
go build -o ./bin/dotenv-sync ./cmd/dotenv-sync
```

On Windows, the binary path becomes `bin\dotenv-sync.exe`.

## 2. Prepare a sample project

Create a sample schema file:

```dotenv
# Application settings
DATABASE_URL=
JWT_SECRET=
PORT=8080
```

Optionally create `.envsync.yaml`:

```yaml
provider: bitwarden
schema_file: .env.example
env_file: .env
mapping:
  DATABASE_URL: shared/dev/database-url
  JWT_SECRET: shared/dev/jwt-secret
```

## 3. Verify prerequisites

```bash
./bin/dotenv-sync doctor
```

Expected result:
- Reports whether the Bitwarden CLI is installed.
- Reports whether the user is logged in and the vault is unlocked.
- Does not print any secret values.

## 4. Preview the sync without writing files

```bash
./bin/dotenv-sync sync --dry-run
```

Expected result:
- Shows only real additions or updates.
- Marks unresolved keys without leaking values.
- Leaves `.env` unchanged.

## 5. Write the local environment file

```bash
./bin/dotenv-sync sync
```

Expected result:
- Creates or updates `.env`.
- Copies safe defaults from `.env.example`.
- Resolves blank schema entries from Bitwarden.
- Preserves comments, ordering, and existing line endings when possible.

## 6. Inspect drift and unresolved values

```bash
./bin/dotenv-sync diff
./bin/dotenv-sync validate
./bin/dotenv-sync missing
```

Expected result:
- `diff` previews differences across schema, local env, and resolved provider values.
- `validate` exits non-zero when drift, malformed files, or missing values exist.
- `missing` lists unresolved keys in a redacted, CI-friendly format.

## 7. Bootstrap or maintain the schema

```bash
./bin/dotenv-sync init
./bin/dotenv-sync reverse --dry-run
./bin/dotenv-sync reverse
```

Expected result:
- `init` creates `.env.example` from `.env` while blanking secret values.
- `reverse --dry-run` previews schema additions only.
- `reverse` writes new blank placeholders back to `.env.example` without
  copying secrets.

## 8. Run the validation suite

```bash
go test ./...
go test ./... -run TestContract
go test ./... -bench .
```

Expected result:
- Unit tests verify parsing, writing, provider, and orchestration behavior.
- Contract tests verify command semantics, exit codes, and output vocabulary.
- Benchmarks confirm the sync, diff, and validate paths meet the documented
  performance budgets.
