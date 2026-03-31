# Provider Contract: Bitwarden write-back for `note_json`

## Scope

- **Provider**: Bitwarden via the `rbw` CLI
- **Item shape**: One repo-scoped Bitwarden login entry per repository
- **Readable modes**:
  - `fields` (existing default, write-capable only for `password`-field aliases)
  - `note_json` (write-capable full-env mode)

## `note_json` Payload Format

The repo-scoped Bitwarden item notes must store a deterministic JSON envelope:

```json
{
  "format": "dotenv-sync/note-json@v1",
  "env": {
    "DATABASE_URL": "postgres://...",
    "JWT_SECRET": "..."
  }
}
```

## Payload Rules

- The `format` field is required and versioned.
- The `env` object contains exact dotenv values keyed by env var name.
- Keys are serialized in deterministic lexical order.
- Payload comparisons are semantic; whitespace-only note edits must not force a
  provider rewrite if the env map is unchanged.

## Provider Call Budget

- Dry-run:
  - at most one provider read of the current payload
- Real push:
  - at most one provider read
  - one `rbw add` or `rbw edit` mutation when required
  - one final `rbw sync`

## `fields` Push Rules

- `ds push` may update field-mode repos only when each pushed schema key maps to
  Bitwarden's built-in `password` field.
- Field-mode push uploads only provider-managed keys currently present in `.env`.
- Multiple env keys may point at the shared `password` field only when their
  local values are identical.
- Custom Bitwarden field mappings remain readable through `ds sync`, but they
  are not writable through `rbw` and must fail with actionable guidance.

## Mutation Rules

- If the repo item is missing, the adapter creates it.
- If the repo item exists and the payload differs semantically, the adapter
  updates it.
- If the payload is already current, the adapter must not mutate the provider.
- All mutation paths must remain scriptable and non-interactive from the `ds`
  operator perspective.

## Read Compatibility Rules

- `ds sync`, `ds diff`, `ds validate`, and `ds missing` must resolve values from
  the `note_json` payload when that mode is configured.
- Repositories that remain on `fields` mode continue to use
  `rbw get <item> --field <key>` semantics unchanged.

## Failure Rules

- Malformed JSON, unsupported `format` values, or incompatible item content must
  produce actionable provider errors.
- Provider errors must never include raw secret values.
- Existing `fields` mode remains readable, but unsupported field-mode writes
  must fail instead of partially mutating the repo item.
