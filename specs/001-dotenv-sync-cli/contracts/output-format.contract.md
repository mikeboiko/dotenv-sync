# Output Format Contract

## Shared Status Vocabulary

All commands use the same normalized status vocabulary:

- `CHECKED`: A prerequisite or input was inspected
- `RESOLVED`: A provider-backed value was found
- `UNCHANGED`: No mutation is required
- `WRITTEN`: A file mutation completed successfully
- `MISSING`: A required secret or prerequisite is unavailable
- `ERROR`: The command cannot continue safely

## Secret Safety Rules

- Raw secret values never appear in stdout or stderr.
- Previews, diffs, diagnostics, and errors use redacted markers such as
  `[RESOLVED]`, `[MISSING]`, or `[STATIC]`.
- Error messages may name a key but must not reveal the secret content.

## Success Message Shape

Success output communicates three elements in a predictable order:

1. The command phase (`CHECKED`, `WRITTEN`, or `UNCHANGED`)
2. The target (`.env`, `.env.example`, provider, config)
3. A compact summary of counts or next steps

Example:
- `WRITTEN .env (added: 2, updated: 1, unchanged: 3)`
- `UNCHANGED .env (already up to date)`

## Dry-Run and Diff Message Shape

Preview output lists only real mutations and unresolved keys.

Example:
- `ADD DATABASE_URL [RESOLVED]`
- `UPDATE JWT_SECRET [RESOLVED]`
- `UNCHANGED PORT [STATIC]`
- `MISSING API_KEY`

## Error Message Shape

Every actionable failure includes:

- `Problem`: what failed
- `Impact`: what the user cannot do yet
- `Action`: how to recover

Example:
- `Problem: Bitwarden database is locked`
- `Impact: sync cannot resolve provider-managed schema keys`
- `Action: run 'rbw unlock' and retry`
