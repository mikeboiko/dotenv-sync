# Provider Contract: Bitwarden write-back for `note_json`

## Scope

- **Provider**: Bitwarden via the `rbw` CLI
- **Item shape**: One repo-scoped Bitwarden login entry per repository
- **Readable modes**:
  - `fields` (existing default, read-only for `ds push`)
  - `note_json` (new write-capable mode)

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
- Existing `fields` mode must be treated as incompatible with `ds push`, not as
  an implicit migration target.
