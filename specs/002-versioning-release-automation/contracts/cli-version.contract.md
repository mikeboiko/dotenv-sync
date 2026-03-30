# CLI Contract: Version inspection

## Commands

### `ds --version`

- **Purpose**: Print a concise, script-friendly version string for the installed
  binary.
- **Exit code**: `0` on success
- **Stdout format**:

```text
ds v0.4.0
```

- **Development-build example**:

```text
ds dev
```

- **Notes**:
  - Must not touch `.env`, `.env.example`, or provider state
  - Must not write to stderr during successful execution

### `ds version`

- **Purpose**: Print detailed build metadata for debugging and release
  verification.
- **Exit code**: `0` on success
- **Stdout format**:

```text
Version: v0.4.0
Commit: 1a2b3c4d
Built: 2026-03-30T11:30:00Z
Platform: linux/amd64
```

- **Development-build example**:

```text
Version: dev
Commit: none
Built: unknown
Platform: darwin/arm64
```

## Shared Rules

- `ds --version` and `ds version` must read from the same embedded metadata
  source.
- The detailed command may expose commit and build-time metadata, but neither
  output path may expose secrets or environment values.
- Extra positional arguments to `ds version` should be treated as usage errors
  consistent with the repository's existing Cobra command behavior.
