# CLI Commands Contract

## Global Rules

- Default executable name: `ds`
- Product and repository name remain `dotenv-sync`
- Default schema file: `.env.example`
- Default local env file: `.env`
- Default config file: `.envsync.yaml`
- Bitwarden access uses the `rbw` CLI in the MVP.
- Standard output carries normal status and preview information.
- Standard error carries actionable failures.
- Secret values are never printed; command output uses redacted or status-only
  representations.
- Process exit codes:
  - `0`: Success or no-op success
  - `1`: Operational failure (configuration, provider, parse, or I/O issue)
  - `2`: Validation or drift failure intended for CI consumption

## sync

**Purpose**: Create or update `.env` from `.env.example`, optional config,
and provider-backed secret resolution.

**Inputs**:

- `--dry-run` to preview without writing
- `--config PATH` to override `.envsync.yaml`
- `--schema PATH` and `--env PATH` for explicit file overrides

**Behavior**:

- Reads `.env.example` as the schema contract.
- Copies explicit defaults directly into `.env`.
- Resolves blank schema values from the provider.
- Preserves comments, ordering, and line endings whenever possible.

**Success outputs**:

- `SYNC CHECKED` for dry-run previews only
- `SYNC WRITTEN`
- `SYNC UNCHANGED`

## diff

**Purpose**: Show redacted differences between schema, local env, and
resolved provider values without writing files.

**Behavior**:

- Highlights adds, updates, unchanged entries, extras, and unresolved values.
- Shows only real mutations.
- Keeps secret values redacted.

## validate

**Purpose**: Verify schema integrity, local env drift, parse validity, and
provider resolution readiness.

**Behavior**:

- Returns exit code `0` when no blocking issues exist.
- Returns exit code `2` when drift, malformed input, duplicates, or
  unresolved secrets are detected.
- Returns exit code `1` for provider or I/O failures that prevent validation.

## doctor

**Purpose**: Diagnose provider and local-environment prerequisites before
sync.

**Behavior**:

- Checks config readability.
- Checks `rbw` presence.
- Checks login and unlock state.
- Reports problem, impact, and next action for each failure.

## init

**Purpose**: Generate `.env.example` from an existing `.env`.

**Behavior**:

- Copies safe explicit defaults into the schema.
- Blanks secret-looking values before writing.
- Preserves user comments and ordering when source formatting allows.

## missing

**Purpose**: Report schema keys that cannot currently be resolved from the
provider.

**Behavior**:

- Outputs the unresolved keys only.
- Returns exit code `2` when any key is unresolved.
- Does not print resolved secret content.

## reverse

**Purpose**: Add new keys from `.env` back into `.env.example` as blank
placeholders.

**Inputs**:

- `--dry-run` to preview schema changes before writing

**Behavior**:

- Adds missing schema keys in a deterministic order.
- Writes blank placeholders only.
- Preserves existing comments and ordering wherever possible.
