# Error Codes Contract

## Process Exit Codes

| Exit Code | Meaning |
|-----------|---------|
| 0 | Command completed successfully, including no-op success |
| 1 | Operational failure prevented the command from completing |
| 2 | Validation or drift issue detected for user or CI action |

## Operator-Visible Error Codes

| Code | Condition | Recovery |
|------|-----------|----------|
| E001 | Provider CLI not installed | Install the provider CLI and rerun `doctor` |
| E002 | Provider not authenticated | Sign in to the provider and rerun the command |
| E003 | Vault locked or unavailable | Unlock the vault and retry |
| E004 | Schema file missing | Create or point to `.env.example`, or use `init` |
| E005 | Secret not found for schema key | Add the secret or mapping, then rerun |
| E006 | Malformed `.env` or `.env.example` | Fix the file formatting and rerun `validate` |
| E007 | Config file invalid | Correct `.envsync.yaml` and rerun `doctor` |
| E008 | Duplicate schema key detected | Remove the duplicate and rerun `validate` |

## Reporting Rules

- Every blocking failure must include one error code.
- Error code text may name a schema key or file path but must not print a secret.
- Validation failures intended for CI should combine exit code `2` with the most
  relevant operator-visible error code in the report body.
