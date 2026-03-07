# Env File Format Contract

## Canonical File Roles

- `.env.example` is the canonical schema file.
- `.env` is the local resolved file for developer tooling.
- Blank schema values indicate provider-managed secrets.
- Explicit schema values indicate safe defaults that may be copied as written.

## Preserved Elements

The implementation MUST preserve the following whenever the source format allows:

- Comment lines
- Blank lines
- Assignment ordering
- Original line ending style (`LF` or `CRLF`)
- Unchanged assignment formatting for no-op runs

## Parsing Rules

- Support comment lines beginning with `#`.
- Support assignment lines in `KEY=VALUE` form.
- Preserve quoted values and surrounding whitespace semantics.
- Flag malformed or duplicate assignments during validation.
- Treat unknown line shapes as parse errors rather than silently rewriting them.

## Writing Rules

- Writes MUST be atomic.
- No-op sync and reverse operations MUST leave file contents byte-for-byte
  unchanged.
- Reverse and init flows MUST write blank placeholders to `.env.example` for
  secret-bearing keys.
- Schema writes MUST never persist provider-resolved secret values.

## Merge Rules

- Sync writes `.env` from schema order first, then applies resolved values.
- Reverse writes only keys present in `.env` and absent from `.env.example`.
- Diff and dry-run output MUST describe the same changes that an actual write
  would perform.
