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

## Or keep `fields` mode for shared aliases

If multiple repos should share one Bitwarden value under different env var names,
keep the default `fields` mode and map each repo's key to Bitwarden's built-in
`password` field:

```yaml
# repo1/.envsync.yaml
item_name: shared-dev
mapping:
  DB_PASSWD: password
```

```yaml
# repo2/.envsync.yaml
item_name: shared-dev
mapping:
  PSWD: password
```

With that setup, `ds push` updates the shared Bitwarden password value and
`ds sync` in either repo reads it back through each repo's local key name.

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
- In `note_json`, `ds` stores the current env map in the repo-scoped Bitwarden
  item notes
- In `fields`, `ds` updates only provider-managed keys currently present in
  `.env` and mapped to Bitwarden's `password` field
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

## Unsupported field-mode guidance

`rbw` cannot script writes to arbitrary custom Bitwarden fields. If a repo stays
on `fields` mode, `ds push` supports only keys mapped to the built-in
`password` field. Repos that need full-env writeback or custom field writes
should switch to `storage_mode: note_json`.
