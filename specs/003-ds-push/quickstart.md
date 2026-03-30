# Quickstart: Bitwarden write-back with `ds push`

## Opt a repository into `note_json` storage

`ds push` needs a write-capable Bitwarden storage mode. In the target repo,
configure `.envsync.yaml` like this:

```yaml
provider: bitwarden
storage_mode: note_json
item_name: Jesse
schema_file: .env.example
env_file: .env
```

If `item_name` is omitted, `ds` still derives it from the Git repository root.
`storage_mode` defaults to `fields`, so `note_json` must be set explicitly.

## Preview what would be uploaded

```bash
cd ~/git/Jesse
ds push --dry-run
```

Expected behavior:

- The command reads `.env`
- It compares the current provider payload, if any
- It reports only adds, updates, unchanged keys, and extras
- It never prints raw values

## Upload the current `.env`

```bash
cd ~/git/Jesse
ds push
```

Expected behavior:

- If the repo item does not exist, `ds` creates it
- `ds` stores the current env map in the repo-scoped Bitwarden item notes
- The provider cache is synced after a successful write
- A second identical run reports `UNCHANGED`
- Extra `.env` keys outside `.env.example` are surfaced as `EXTRA` without
  mutating the schema

## Rebuild `.env` from the pushed provider data

```bash
cd ~/git/Jesse
rm .env
ds sync
```

With `storage_mode: note_json`, `ds sync` resolves provider-managed values from
the same repo-scoped Bitwarden payload that `ds push` writes.

## Incompatible-mode guidance

If the repo still uses the default field-based Bitwarden layout, `ds push` must
fail with actionable guidance rather than silently writing a second, unread
secret representation. The recovery path is to opt the repo into
`storage_mode: note_json` first.
