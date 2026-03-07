# Research: dotenv-sync CLI MVP

## Decision 1: Implement the MVP as a cross-platform Go CLI

- **Decision**: Use Go 1.22 to ship a single CLI binary for Linux, macOS,
  and Windows.
- **Rationale**: The user explicitly requested a cross-platform Go
  application, and Go provides predictable cross-compilation, fast startup,
  easy subprocess control, and a small operational footprint for local
  developer tooling.
- **Alternatives considered**:
  - Python CLI: Rejected because it adds runtime packaging and interpreter
    dependencies that weaken the zero-friction onboarding goal.
  - Rust CLI: Rejected because it would satisfy the technical needs but does
    not match the requested language and would slow initial iteration for the
    MVP.

## Decision 2: Use Cobra for command routing and help output

- **Decision**: Use `github.com/spf13/cobra` for the root command and
  subcommand tree.
- **Rationale**: Cobra gives stable subcommand composition, built-in help
  text, shell completion support, and consistent flag parsing for sync,
  diff, validate, doctor, init, missing, and reverse flows without custom
  routing logic.
- **Alternatives considered**:
  - Standard `flag` package: Rejected because multi-command UX and help
    handling would become repetitive and harder to maintain.
  - `urfave/cli`: Rejected because Cobra is the more common choice for Go
    CLIs with deep subcommand trees and strong documentation expectations.

## Decision 3: Build a custom envfile tokenizer and writer

- **Decision**: Parse `.env` and `.env.example` into ordered line tokens and
  use a custom writer for merges and rewrites.
- **Rationale**: The constitution requires comment, ordering, no-op, and
  line-ending preservation, which common map-based dotenv libraries do not
  guarantee during rewrites. A custom token model lets the implementation
  keep formatting stable while still resolving values and generating diffs.
- **Alternatives considered**:
  - End-to-end use of a dotenv library: Rejected because most libraries
    discard comments and reorder keys during writes.
  - In-place string replacement: Rejected because it is fragile around
    quoting, duplicates, multiline values, and reverse-sync behavior.

## Decision 4: Integrate Bitwarden through the `rbw` CLI behind a provider interface

- **Decision**: Use `exec.CommandContext` to invoke the `rbw` CLI from a
  `provider.Provider` adapter rather than calling shell wrappers or
  implementing a direct Bitwarden API client for the MVP.
- **Rationale**: The installed `rbw` CLI already exposes the login, unlock,
  sync, list, and get operations needed for the MVP. Wrapping it in a
  provider interface keeps Bitwarden-specific logic out of the sync engine
  and preserves an upgrade path for future providers or future Bitwarden CLI
  variants.
- **Alternatives considered**:
  - Direct Bitwarden API integration: Rejected because it would duplicate
    auth and session complexity that the CLI already solves.
  - The official `bw` CLI in the MVP: Deferred to the roadmap so the first
    release can align with the user's installed tooling without expanding the
    MVP surface.
  - Shell-script glue: Rejected because it is harder to make deterministic,
    secure, and cross-platform.

## Decision 5: Use `.envsync.yaml` for optional project configuration

- **Decision**: Support an optional `.envsync.yaml` file for provider
  selection, vault hints, schema and env file paths, and explicit key
  mapping.
- **Rationale**: The user already described a YAML mapping file, and YAML is
  readable for teams while remaining optional for zero-config onboarding. The
  dependency cost of `gopkg.in/yaml.v3` is justified because it preserves the
  expected project-level configuration UX.
- **Alternatives considered**:
  - JSON config: Rejected because it does not match the requested
    configuration style and is less friendly for hand editing.
  - Flags-only configuration: Rejected because persistent key mapping
    becomes cumbersome and hard to share across a team.

## Decision 6: Preserve platform-specific file and process behavior explicitly

- **Decision**: Use `filepath` for paths, preserve detected line endings,
  perform atomic rewrites with temp files and rename, and avoid shell
  execution entirely.
- **Rationale**: Windows and Unix-like environments differ in separators,
  default shells, and line endings, but the CLI must behave the same across
  all three target platforms. Explicit path handling and atomic writes reduce
  cross-platform bugs and improve confidence in deterministic rewrites.
- **Alternatives considered**:
  - Always writing LF line endings: Rejected because it would create
    unnecessary diffs on Windows-heavy projects.
  - Executing provider calls through `sh` or `powershell`: Rejected because
    it introduces quoting risk and inconsistent behavior across platforms.

## Decision 7: Use failing-first Go tests plus CLI contract and integration suites

- **Decision**: Build the test strategy around `go test`, package-level unit
  tests, CLI contract tests, integration tests with fixture projects, and
  benchmark coverage for hot paths.
- **Rationale**: The constitution requires failing-first automation and
  explicit negative-path coverage for provider failures and redaction. Go's
  standard test tooling is enough to support unit, integration, and benchmark
  workflows while keeping dependencies minimal.
- **Alternatives considered**:
  - Manual verification only: Rejected because it cannot reliably prove
    deterministic rewrites or secret-safe output.
  - Unit tests without CLI contracts: Rejected because command semantics and
    exit-code behavior are part of the public surface of the tool.

## Decision 8: Ship `ds` as the default executable name

- **Decision**: Use `ds` as the default binary name while retaining
  `dotenv-sync` as the product and repository name.
- **Rationale**: The user explicitly called out typing ergonomics, and the
  shorter executable improves the everyday developer experience without
  changing the conceptual identity of the project.
- **Alternatives considered**:
  - Shipping only a `dotenv-sync` binary: Rejected because it adds avoidable
    friction to the main command path.
  - Renaming the whole project to `ds`: Rejected because it would discard the
    more descriptive product name.

## Decision 9: Keep `bw` CLI compatibility in the roadmap

- **Decision**: Treat support for the official `bw` CLI as a post-MVP
  Bitwarden compatibility enhancement.
- **Rationale**: The provider interface and package structure can support an
  alternate Bitwarden client later, but the MVP remains focused on the
  user's existing `rbw` workflow and the smallest high-value integration.
- **Alternatives considered**:
  - Shipping both `rbw` and `bw` support in MVP: Rejected because it expands
    testing, diagnostics, and UX complexity before the core workflow is
    proven.
