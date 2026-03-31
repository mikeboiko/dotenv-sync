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
- Bitwarden `storage_mode` must be either:
  - `note_json`, or
  - `fields` with pushed keys mapped to Bitwarden's built-in `password` field

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
- Field-mode repos may push shared aliases through the Bitwarden `password`
  field, while unsupported custom-field mappings must fail with actionable
  guidance.
- Summary lines use the shared status vocabulary and target the provider item,
  for example:

```text
CHECKED bitwarden:Jesse (added: 2, updated: 1, extra: 1, dry-run)
WRITTEN bitwarden:Jesse (added: 2, updated: 1, extra: 1)
UNCHANGED bitwarden:Jesse (already up to date)
```

## Failure Rules

- If the repo stays on `fields` mode but maps pushed keys to custom Bitwarden
  fields, `ds push` must fail with actionable guidance to use `password` or
  switch to `note_json`.
- If `.env` is missing, malformed, or contains duplicate keys, `ds push` must
  fail before any provider mutation.
- If the provider payload is malformed, `ds push` must fail with a repair action
  rather than overwriting it blindly.
- Successful execution must not write to stderr.
