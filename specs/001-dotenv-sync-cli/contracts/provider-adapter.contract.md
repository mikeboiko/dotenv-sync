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
- `CheckReadiness(...)` verifies CLI installation, authentication, and unlock state
- `Resolve(...)` returns one resolution result for a schema key or mapped lookup
- `ResolveMany(...)` may optimize repeated lookups while preserving one logical
  lookup per distinct key per command

## Result Guarantees

- Readiness failures must map to stable error codes.
- Resolution failures must differentiate between `missing`, `unmapped`, and
  provider errors.
- Returned errors must never include raw secret values.
- Adapters must be safe to replace with mocks for tests.

## Bitwarden-Specific Expectations

- Readiness checks use the `bw` CLI directly.
- Mapping may point a schema key to a custom Bitwarden item or field reference.
- The adapter caches repeated lookups within a command to satisfy the
  performance budget.
- Unlock and login failures surface actionable recovery guidance.
