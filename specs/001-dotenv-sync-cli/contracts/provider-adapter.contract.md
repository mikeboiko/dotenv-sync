# Provider Adapter Contract

## Responsibilities

A provider adapter is responsible for:

- Reporting readiness and prerequisite failures
- Resolving schema keys into secret values or missing-value results
- Hiding provider-specific transport details from the sync engine
- Returning redaction-safe diagnostics

## Core Interface Expectations

The initial Go interface is expected to support these behaviors:

- `Name()` returns the provider name
- `CheckReadiness(...)` verifies CLI installation and provider-specific access
  state such as login, unlock, database, or group readiness
- `Resolve(...)` returns one resolution result for a schema key and provider
  lookup reference
- `ResolveMany(...)` may optimize repeated lookups while preserving one logical
  lookup per distinct key per command

## Result Guarantees

- Readiness failures must map to stable error codes.
- Resolution failures must differentiate between `missing`, `unmapped`, and
  provider errors.
- Returned errors must never include raw secret values.
- Adapters must be safe to replace with mocks for tests.

## Bitwarden-Specific Expectations

- Readiness checks use the `rbw` CLI directly.
- The default Bitwarden item name is derived from the Git repository root
  directory name, falling back to the current working-directory name when Git
  metadata is unavailable.
- Unmapped schema keys resolve as `rbw get <item-name> --field <schema-key>`.
- Config may override the Bitwarden item name and may remap individual schema
  keys to alternate field names within that item.
- The adapter caches repeated lookups within a command to satisfy the
  performance budget.
- Unlock and login failures surface actionable recovery guidance.

## KeePass-Specific Expectations

- Readiness checks use the `keepassxc-cli` binary directly.
- The configured database path and group define the lookup root for provider
  values.
- Unmapped schema keys resolve as `<group>/<schema-key>` entry paths.
- KeePass lookup paths are derived from `<group>/<schema-key>`; current
  Bitwarden-style `mapping` overrides do not change KeePass entry paths.
- The adapter prompts for the KeePass master password at most once per command
  and reuses it for later lookups in the same process.
- The adapter caches repeated lookups within a command to satisfy the
  performance budget.

## Roadmap Note

- Future roadmap work may add new providers, or alternate clients for existing
  providers, behind the same provider interface without changing sync-engine
  contracts.
