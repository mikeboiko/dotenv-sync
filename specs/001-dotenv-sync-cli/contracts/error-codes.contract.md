# Error Codes Contract

## Process Exit Codes

| Exit Code | Meaning                                                   |
| --------- | --------------------------------------------------------- |
| 0         | Command completed successfully, including no-op success   |
| 1         | Operational failure prevented the command from completing |
| 2         | Validation or drift issue detected for user or CI action  |

## Operator-Visible Error Codes

Error codes are stable within a command surface. Some codes are reused across
related flows, so the exact message and recovery guidance depend on the command
and provider context.

| Code | Surface                    | Current meaning                                                                       | Recovery                                                                |
| ---- | -------------------------- | ------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| E001 | provider readiness         | Required provider CLI missing (`rbw` or `keepassxc-cli`)                              | Install or expose the configured provider CLI on `PATH`, then rerun     |
| E002 | provider readiness         | Provider access prerequisite missing (for example Bitwarden login or KeePass DB path) | Restore the missing provider prerequisite and rerun                     |
| E003 | provider readiness / I/O   | Provider locked, unreadable, or provider read/write operation failed                  | Unlock or repair the configured provider and retry                      |
| E004 | schema / provider location | Schema file missing, or configured provider location/group unavailable                | Create the schema, or fix the provider path/group and rerun             |
| E005 | resolution                 | Secret not found for schema key                                                       | Add the secret or adjust the lookup mapping, then rerun                 |
| E006 | file parsing / file I/O    | Local env, schema, or config file could not be read, parsed, or written               | Fix the file contents or permissions and rerun                          |
| E007 | config / capability gating | Config file invalid, or a command requires a different provider capability            | Correct `.envsync.yaml`, or use a compatible provider configuration     |
| E008 | validation / capability    | Duplicate schema key detected, or `ds push` used with an unsupported provider         | Remove the duplicate, or switch to a provider that supports the command |
| E009 | push storage mode          | Bitwarden push storage mode incompatible with the requested write flow                | Use a compatible Bitwarden storage mode and rerun                       |
| E010 | Bitwarden note-json        | Repo-scoped Bitwarden `note_json` payload malformed or could not be serialized        | Repair or recreate the Bitwarden item notes, then retry                 |
| E011 | Bitwarden fields-mode push | `fields`-mode push cannot safely target custom or conflicting password mappings       | Map pushed keys to `password`, or switch to `storage_mode: note_json`   |

## Reporting Rules

- Every blocking failure must include one error code.
- Error code text may name a schema key or file path but must not print a secret.
- Validation failures intended for CI should combine exit code `2` with the most
  relevant operator-visible error code in the report body.
