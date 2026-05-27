# Data Model: dotenv-sync CLI MVP

## Overview

The MVP operates on ordered envfile documents, provider configuration,
resolution results, and command outcomes. These entities are in-memory domain
models that support sync, diff, validate, doctor, init, missing, and reverse
flows.

## Entities

### EnvironmentDocument

Represents a parsed `.env`-style file.

| Field          | Type                     | Description                                          |
| -------------- | ------------------------ | ---------------------------------------------------- |
| path           | string                   | Absolute or project-relative path to the source file |
| kind           | enum (`schema`, `local`) | Distinguishes `.env.example` from `.env`             |
| lineEnding     | enum (`LF`, `CRLF`)      | Original line ending style to preserve on write      |
| lines          | []EnvironmentLine        | Ordered list of parsed lines                         |
| hasParseErrors | bool                     | Indicates whether malformed lines were encountered   |

**Validation rules**:

- `path` must reference a single file in the target project.
- `kind` must align with the command context.
- `hasParseErrors` blocks write operations until surfaced to the user.

### EnvironmentLine

Represents one ordered line in a dotenv file.

| Field             | Type                                    | Description                                          |
| ----------------- | --------------------------------------- | ---------------------------------------------------- |
| index             | int                                     | Original position in the file                        |
| raw               | string                                  | Original line content for no-op comparisons          |
| lineType          | enum (`blank`, `comment`, `assignment`) | Parsed line category                                 |
| key               | string                                  | Environment variable name when `lineType=assignment` |
| value             | string                                  | Parsed literal value when present                    |
| inlineComment     | string                                  | Trailing comment content when present                |
| managedByProvider | bool                                    | True when a schema assignment is intentionally blank |

**Validation rules**:

- Assignment keys must match the expected environment variable naming rules.
- Duplicate assignment keys are reported as validation issues.
- `managedByProvider` may only be true for schema lines.

### AppConfig

Represents optional project configuration from `.envsync.yaml`.

| Field           | Type              | Description                                                                               |
| --------------- | ----------------- | ----------------------------------------------------------------------------------------- |
| provider        | string            | Selected provider family, defaulting to `bitwarden`                                       |
| schemaFile      | string            | Override path for `.env.example`                                                          |
| envFile         | string            | Override path for `.env`                                                                  |
| itemName        | string            | Override for the default repo-scoped Bitwarden item name used for provider lookups        |
| storageMode     | string            | Bitwarden storage strategy (`fields` or `note_json`) when provider write-back is used     |
| vault           | string            | Reserved provider-specific scope hint; current built-in adapters do not use it            |
| mapping         | map[string]string | Schema key to provider-specific lookup override for adapters that support alternate names |
| keepassDatabase | string            | Path to the KeePass `.kdbx` file when `provider=keepass`                                  |
| keepassGroup    | string            | Group inside the KeePass database that stores env-var entries                             |

**Validation rules**:

- `provider` must map to a registered adapter.
- `schemaFile` and `envFile` must not point to the same path.
- `itemName` defaults to the Git repository root directory name, or the current
  working-directory name when Git metadata is unavailable, when the selected
  provider uses repo-scoped item lookups.
- `storageMode` applies only to Bitwarden-backed write flows and defaults to
  `fields`.
- Mapping keys must be unique and reference valid schema keys when used.
- Mapping values override the default lookup name for a schema key within the
  selected provider's addressing model; current built-in adapters use mapping
  for Bitwarden field-name overrides only.
- `vault` is reserved for future provider-specific scoping and is currently
  ignored by the built-in Bitwarden and KeePass adapters.
- `keepassDatabase` and `keepassGroup` are required only when
  `provider=keepass`.

### ProviderStatus

Represents current readiness of the configured secret provider.

| Field         | Type   | Description                                                        |
| ------------- | ------ | ------------------------------------------------------------------ |
| provider      | string | Provider name under test                                           |
| cliInstalled  | bool   | Whether the configured provider CLI or local integration is usable |
| authenticated | bool   | Whether the provider session or equivalent auth state is ready     |
| unlocked      | bool   | Whether the provider data source is unlocked and readable          |
| message       | string | Human-readable diagnostic summary                                  |

**State transitions**:

- `unknown -> unavailable` when the configured provider CLI or database cannot
  be found.
- `unknown -> installed` when the CLI is present.
- `installed -> authenticated` when provider login or equivalent readiness is
  confirmed.
- `authenticated -> unlocked` when secrets can be queried.

### SecretResolution

Represents the result of resolving one schema key.

| Field       | Type                                                        | Description                                          |
| ----------- | ----------------------------------------------------------- | ---------------------------------------------------- |
| key         | string                                                      | Schema key being resolved                            |
| source      | enum (`static`, `provider`, `missing`, `unmapped`, `error`) | Resolution outcome                                   |
| providerRef | string                                                      | Lookup name or entry path used for the provider call |
| value       | string                                                      | Resolved secret or copied static literal             |
| redacted    | string                                                      | Safe display form for logs and previews              |
| issueCode   | string                                                      | Optional error or warning code                       |

**Validation rules**:

- `redacted` must never expose `value` when `source=provider`.
- `source=static` requires the value to come from the schema file.
- `source=missing`, `unmapped`, and `error` must populate `issueCode`.

### SyncPlan

Represents the proposed result of a command before any file write.

| Field          | Type                                                                                 | Description                                        |
| -------------- | ------------------------------------------------------------------------------------ | -------------------------------------------------- |
| mode           | enum (`sync`, `dry-run`, `diff`, `validate`, `doctor`, `init`, `missing`, `reverse`) | Requested command mode                             |
| schema         | EnvironmentDocument                                                                  | Parsed schema input                                |
| localEnv       | EnvironmentDocument                                                                  | Parsed local env input when present                |
| config         | AppConfig                                                                            | Effective configuration                            |
| providerStatus | ProviderStatus                                                                       | Provider readiness state                           |
| resolutions    | []SecretResolution                                                                   | Resolution result per schema key                   |
| changes        | []ChangeRecord                                                                       | Proposed additions, updates, removals, or warnings |
| writeRequired  | bool                                                                                 | Whether a file mutation would occur                |

**Validation rules**:

- `doctor` may omit `schema` and `localEnv` writes but still requires config and provider status.
- `sync` and `reverse` require schema and local env documents when writing.
- `writeRequired=false` with no issues results in a no-op success state.

### ChangeRecord

Represents one user-visible change or validation finding.

| Field      | Type                                                             | Description                     |
| ---------- | ---------------------------------------------------------------- | ------------------------------- |
| key        | string                                                           | Affected environment variable   |
| changeType | enum (`add`, `update`, `unchanged`, `extra`, `missing`, `error`) | Classification used in reports  |
| before     | string                                                           | Previous redacted display value |
| after      | string                                                           | New redacted display value      |
| file       | enum (`schema`, `local`)                                         | File affected by the change     |
| message    | string                                                           | Human-readable explanation      |

### ValidationIssue

Represents a blocking or non-blocking problem found during parsing,
validation, or provider checks.

| Field    | Type                      | Description                           |
| -------- | ------------------------- | ------------------------------------- |
| code     | string                    | Stable issue identifier               |
| severity | enum (`error`, `warning`) | Impact level                          |
| file     | string                    | File path or provider target involved |
| key      | string                    | Related key when applicable           |
| message  | string                    | Problem description                   |
| action   | string                    | Suggested next step for the user      |

## Relationships

- `AppConfig` shapes how `SyncPlan` loads the schema, env file, and provider.
- `EnvironmentDocument` contains ordered `EnvironmentLine` records.
- `SyncPlan` compares `schema` and `localEnv`, then enriches the result with
  `ProviderStatus`, `SecretResolution`, `ChangeRecord`, and `ValidationIssue`.
- `doctor` focuses primarily on `ProviderStatus` and config validation, while
  `init` and `reverse` transform `EnvironmentDocument` instances.
