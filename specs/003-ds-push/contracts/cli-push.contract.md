# CLI Contract: `ds push`

## Command

### `ds push`

- **Purpose**: Upload the current `.env` values into the repo-scoped Bitwarden
  entry without printing secret content.
- **Inputs**:
  - `--dry-run` to preview provider mutations without writing
  - existing global flags: `--config`, `--schema`, `--env`
- **Exit codes**:
  - `0` on success or no-op success
  - `1` on config, provider, parse, or I/O failures

## Preconditions

- `.env` must exist and parse successfully
- Bitwarden readiness must pass through the existing `rbw` checks
- Bitwarden `storage_mode` must be `note_json`

## Behavior

1. Read `.env` as the upload source.
2. Read `.env.example` as schema context only.
3. Load the current repo-scoped Bitwarden payload for the configured item.
4. Compute adds, updates, unchanged keys, and extras without printing values.
5. If `--dry-run` is set, print the preview and do not mutate the provider.
6. If the provider payload already matches the local env map, print `UNCHANGED`
   and do not mutate the provider.
7. Otherwise create or update the Bitwarden item and refresh the local provider
   cache.

## Stdout Expectations

- Per-key preview lines use change verbs such as `ADD`, `UPDATE`, `UNCHANGED`,
  or `EXTRA`.
- Per-key lines may include redacted markers such as `[REDACTED]`, but must
  never print raw values.
- Field-mode repos remain the default elsewhere in the product, so `ds push`
  must clearly reject them until `storage_mode: note_json` is configured.
- Summary lines use the shared status vocabulary and target the provider item,
  for example:

```text
CHECKED bitwarden:Jesse (added: 2, updated: 1, extra: 1, dry-run)
WRITTEN bitwarden:Jesse (added: 2, updated: 1, extra: 1)
UNCHANGED bitwarden:Jesse (already up to date)
```

## Failure Rules

- If the repo is still on `fields` mode, `ds push` must fail with actionable
  guidance to enable `note_json`.
- If `.env` is missing, malformed, or contains duplicate keys, `ds push` must
  fail before any provider mutation.
- If the provider payload is malformed, `ds push` must fail with a repair action
  rather than overwriting it blindly.
- Successful execution must not write to stderr.
