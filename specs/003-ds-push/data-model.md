# Data Model: Bitwarden write-back with `ds push`

## Push Storage Mode

- **Purpose**: Captures how a repository stores Bitwarden-backed values.
- **Fields**:
  - `name`: enum value `fields` or `note_json`
  - `item_name`: repo-scoped Bitwarden item name derived from config or repo root
  - `config_source`: whether the mode came from defaults or `.envsync.yaml`
- **Validation rules**:
  - `fields` remains the default and existing read behavior
  - `note_json` is required for `ds push`
  - unsupported values fail config loading before command execution

## Local Env Snapshot

- **Purpose**: Represents the current `.env` contents used as the upload source.
- **Fields**:
  - `path`: absolute file path
  - `entries`: ordered list or map of key/value pairs parsed from `.env`
  - `duplicates`: duplicate keys detected during parsing
  - `parse_errors`: formatting errors that block push
  - `extra_keys`: keys present in `.env` but absent from `.env.example`
- **Validation rules**:
  - parse errors and duplicates block write operations
  - values are treated as opaque strings and must round-trip exactly

## Provider Env Payload

- **Purpose**: Represents the repo-scoped Bitwarden state relevant to push and
  note-backed resolution.
- **Fields**:
  - `item_name`: Bitwarden item name
  - `exists`: whether the Bitwarden item already exists
  - `storage_mode`: `fields` or `note_json`
  - `raw_notes`: raw provider notes text when `note_json` is configured
  - `format`: payload format identifier such as `dotenv-sync/note-json@v1`
  - `env`: parsed env map from the provider payload
  - `issue`: malformed, missing, or incompatible payload classification
- **Validation rules**:
  - `fields` mode has no provider payload map for push and is incompatible with
    upload
  - `note_json` payloads must parse as JSON and contain a supported `format`
  - malformed payloads block sync and push until repaired or replaced

## Push Plan

- **Purpose**: Represents the computed result of `ds push` before any provider
  mutation occurs.
- **Fields**:
  - `mode`: command mode such as `push`
  - `dry_run`: whether this is preview-only
  - `source`: associated `Local Env Snapshot`
  - `provider`: associated `Provider Env Payload`
  - `changes`: ordered list of `Push Change` records
  - `write_required`: whether the provider needs to be created or updated
  - `create_required`: whether the repo item does not exist yet
  - `issues`: blocking validation or provider issues
- **Validation rules**:
  - `write_required=false` with no issues results in `UNCHANGED`
  - dry-run and real-run must compute the same `changes`
  - blocking issues prevent provider mutation

## Push Change

- **Purpose**: Describes one previewable mutation or status line for a single
  env key.
- **Fields**:
  - `key`: env key name
  - `change_type`: enum `add`, `update`, `unchanged`, `extra`, or `error`
  - `before_marker`: redacted marker for prior provider state
  - `after_marker`: redacted marker for target provider state
  - `schema_member`: whether the key exists in `.env.example`
  - `message`: optional operator guidance
- **Validation rules**:
  - marker fields must never contain raw values
  - `extra` keys remain pushable but must stay visible in previews

## Relationships

- One **Push Storage Mode** governs one repo-scoped **Provider Env Payload**.
- One **Local Env Snapshot** and one **Provider Env Payload** produce one
  **Push Plan**.
- One **Push Plan** contains many **Push Changes**.
- In `note_json` mode, one **Provider Env Payload** can later feed the existing
  read flows (`sync`, `diff`, `validate`, `missing`).

## State Transitions

- **Provider Env Payload**: `missing -> created -> current -> updated`
- **Push Plan**: `planned -> validated -> written`
- **Push Plan** failure path: `planned|validated -> failed`
