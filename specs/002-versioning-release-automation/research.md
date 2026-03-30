# Research: Version reporting and release automation

## Decision 1: Inject version metadata at build time with mutable Go variables

- **Decision**: Replace the current constant-only version placeholder with
  mutable package variables for semantic version, commit SHA, and build time, and
  populate them through `go build -ldflags -X`.
- **Rationale**: This is the standard Go pattern for CLI version reporting,
  avoids per-release source edits, and keeps local development builds usable by
  falling back to deterministic defaults like `dev`, `none`, and `unknown`.
- **Alternatives considered**:
  - Hard-code versions in source on each release: rejected because it is easy to
    forget and creates noisy commits.
  - Read Git metadata at runtime: rejected because installed binaries may not run
    in a Git checkout and version inspection should not depend on repository
    state.

## Decision 2: Support both a quick flag and a detailed command

- **Decision**: Provide a concise `ds --version` output path and a richer
  `ds version` command that prints detailed build metadata from the same shared
  source.
- **Rationale**: The quick flag matches common CLI expectations for scripts and
  support checks, while the subcommand provides enough detail for debugging and
  release verification without inventing multiple metadata sources.
- **Alternatives considered**:
  - Expose only `ds version`: rejected because users explicitly asked for
    `--version` and many CLIs support the flag.
  - Expose only `--version`: rejected because one-line output is too limited for
    commit and build-time inspection.

## Decision 3: Use semver Git tags as the release source of truth

- **Decision**: Treat the highest reachable semver tag as the current release,
  compute the next version from a `major`, `minor`, or `patch` bump, ignore
  unrelated tags, and use `v0.0.0` as the no-tag baseline.
- **Rationale**: Git tags are the canonical release identifier for a Go CLI and
  avoid splitting version truth between source files, workflows, and GitHub
  Releases.
- **Alternatives considered**:
  - Maintain a checked-in `VERSION` file: rejected because it duplicates Git tag
    state and invites drift.
  - Use Release Please or commit-convention-driven automation: rejected because
    the user chose a manual GitHub Actions bump workflow and the repo does not
    currently enforce commit conventions.

## Decision 4: Keep release automation lean with GitHub Actions, `go build`, and `gh`

- **Decision**: Implement release automation as a `workflow_dispatch` GitHub
  Actions workflow that runs `go test ./...`, builds the supported target
  matrix with injected metadata, and publishes a GitHub release through `gh`.
- **Rationale**: This uses tooling already available on GitHub-hosted runners,
  avoids adding runtime dependencies to `ds`, and keeps the release path readable
  for maintainers.
- **Alternatives considered**:
  - Add GoReleaser: rejected for the first cut because it introduces additional
    configuration and third-party workflow dependencies beyond the current needs.
  - Push tags and releases from a local maintainer machine: rejected because it
    is less repeatable and harder to audit.

## Decision 5: Publish only after validation and builds succeed

- **Decision**: The workflow must compute the next version, run tests, build all
  release artifacts, and only then create the tag and GitHub release.
- **Rationale**: This prevents partial releases where a tag exists but the build
  artifacts or release publication failed, which would make `ds --version` trust
  worse rather than better.
- **Alternatives considered**:
  - Create the tag before builds start: rejected because a failed build would
    leave behind misleading release state.
  - Publish per-platform artifacts incrementally: rejected because the release
    should be all-or-nothing for maintainers and users.

## Decision 6: Document local injection and maintainer release steps

- **Decision**: Update README and quickstart examples so contributors can build a
  local versioned binary, inspect it, and trigger the manual release workflow.
- **Rationale**: Versioning only helps if users and maintainers can discover the
  supported commands and release path without reading implementation details.
- **Alternatives considered**:
  - Document only the GitHub workflow: rejected because contributors also need a
    reproducible local verification path.
