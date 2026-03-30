# Research: Bitwarden write-back with `ds push`

## Decision 1: Add an opt-in `note_json` storage mode instead of changing the existing default

- **Decision**: Keep the current Bitwarden `fields` layout as the default for
  backward compatibility and introduce an opt-in `note_json` storage mode for
  repos that want `rbw`-backed write-back through `ds push`.
- **Rationale**: Existing repositories already rely on repo-scoped field lookups
  such as `rbw get <item> --field <KEY>`. Preserving that default avoids
  breaking current `sync`, `diff`, and `validate` behavior while still enabling
  a write-capable mode for new or migrated repos.
- **Alternatives considered**:
  - **Change the default storage layout to note-backed values for everyone**:
    rejected because it would silently change existing Bitwarden lookups.
  - **Add a `bw` fallback for writes**: rejected because the project explicitly
    standardized on `rbw` and moved `bw` support to the roadmap.

## Decision 2: Store repo env values as deterministic JSON in the Bitwarden item's notes

- **Decision**: In `note_json` mode, store the repo's env values in the
  repo-scoped Bitwarden item's notes as a deterministic JSON envelope such as:
  `{"format":"dotenv-sync/note-json@v1","env":{...}}`.
- **Rationale**: `rbw` already supports reading notes through its field access
  path and its scripted add/edit flows can write password-plus-notes content.
  JSON gives the feature a stable, schema-extensible format and allows semantic
  comparisons without relying on raw string equality.
- **Alternatives considered**:
  - **Custom Bitwarden fields inside the existing repo item**: rejected because
    the current `rbw` CLI surface used by this repo exposes field reads cleanly
    but does not expose an equally scriptable custom-field mutation flow.
  - **Multiple Bitwarden items, one per env var**: rejected because it abandons
    the repo-scoped lookup model and increases provider calls.
  - **Plain dotenv text in notes**: rejected because canonical JSON is easier to
    serialize deterministically and validate safely.

## Decision 3: Treat `.env` as the push source while keeping `.env.example` as schema context

- **Decision**: `ds push` will read `.env` as the upload source, compare it with
  `.env.example`, surface extra keys explicitly, and never mutate either file.
- **Rationale**: The user's requested workflow is "save what is currently in
  `.env` to Bitwarden." At the same time, the project constitution requires
  `.env.example` to remain the schema contract, so push uses the schema for
  context and warnings without rewriting it.
- **Alternatives considered**:
  - **Push only keys present in `.env.example`**: rejected because it would fail
    the user's explicit "all env vars currently in `.env`" expectation.
  - **Update `.env.example` during push**: rejected because write-back to the
    provider should not also mutate local files or schema.

## Decision 4: Automate `rbw add` and `rbw edit` through a temporary editor helper

- **Decision**: Drive `rbw add` and `rbw edit` non-interactively by running them
  with `VISUAL` or `EDITOR` set to a temporary helper script that writes the
  exact password-plus-notes content expected by `rbw`.
- **Rationale**: The local `rbw` CLI uses editor-driven add and edit flows for
  login entries. A temporary helper script keeps `ds` cross-platform, avoids
  shell-specific hacks, and lets the provider adapter remain a wrapper around
  the installed `rbw` executable.
- **Alternatives considered**:
  - **Prompt users to edit manually during `ds push`**: rejected because it
    breaks deterministic automation and makes previews impossible to trust.
  - **Patch the local Bitwarden cache directly**: rejected because it is unsafe
    and bypasses `rbw`'s normal sync/update behavior.

## Decision 5: Add a write-capable provider contract rather than embedding push logic in the CLI

- **Decision**: Introduce a write-capable provider contract for push planning and
  mutation while keeping CLI routing, sync orchestration, and provider-specific
  `rbw` behavior in their existing focused modules.
- **Rationale**: The constitution requires explicit module boundaries. Push needs
  reusable logic for preview, no-op detection, and provider errors, and that
  logic belongs in the sync/provider layers rather than inside Cobra handlers.
- **Alternatives considered**:
  - **Implement push directly inside `internal/cli/push.go`**: rejected because
    it would mix I/O, provider mutation, redaction, and planning logic.
  - **Overload the existing read-only provider interface with CLI-specific
    helpers**: rejected because write semantics are different enough to justify a
    separate capability boundary.

## Decision 6: Keep the provider budget to one read, one mutation, and one sync

- **Decision**: In `note_json` mode, `ds push` should do at most one provider
  fetch of the current notes payload, one add-or-edit mutation if needed, and
  one final `rbw sync` after a successful write.
- **Rationale**: This keeps push aligned with the project's performance budget
  and makes dry-run vs. real-run behavior easy to reason about. A repo-scoped
  note payload also minimizes remote round-trips compared with per-key writes.
- **Alternatives considered**:
  - **One provider write per env var**: rejected because it would be slower,
    harder to preview, and contrary to the repo-scoped storage choice.
  - **Re-read the provider after every mutation**: rejected because a single
    final sync is enough to refresh the local cache for later commands.
