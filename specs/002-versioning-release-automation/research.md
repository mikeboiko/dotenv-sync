# Research: Automatic minor release automation

## Decision 1: Reuse embedded version metadata only for release verification

- **Decision**: Keep build metadata injection and the existing CLI version
  surfaces as the source used to verify built artifacts before publication,
  rather than introducing a second release-verification mechanism.
- **Rationale**: The binary already exposes trustworthy metadata through
  `ds --version` and `ds version`, so the release workflow can validate a built
  artifact without adding runtime dependencies or new operator-facing commands.
- **Alternatives considered**:
  - Add a separate release verification binary or script output format: rejected
    because it duplicates version truth and adds maintenance overhead.
  - Skip artifact self-verification entirely: rejected because automatic
    publication should confirm that released binaries report the version they
    contain.

## Decision 2: Trigger release automation on pushes to the default branch and always bump minor

- **Decision**: Replace manual `workflow_dispatch` release creation with an
  automatic workflow that runs on pushes to `main` and always computes the next
  minor version.
- **Rationale**: The user asked for automatic minor bumps on every push to
  `main`, and eliminating manual bump selection removes the remaining human step
  in the release path.
- **Alternatives considered**:
  - Keep manual dispatch with a default `minor` input: rejected because it still
    requires a maintainer to trigger each release.
  - Use commit messages or labels to select major/minor/patch bumps: rejected
    because the requested behavior is minor-only and the repository does not rely
    on commit-convention parsing today.

## Decision 3: Continue treating semver Git tags as the release source of truth

- **Decision**: Compute the next minor version from the latest reachable semver
  tag, ignore unrelated tags, and use `v0.0.0` when no semver tag exists.
- **Rationale**: Git tags remain the canonical release identifier for a Go CLI,
  and minor-only automation can build on the existing tag-based logic without
  splitting release truth across source files, workflows, and GitHub Releases.
- **Alternatives considered**:
  - Maintain a checked-in `VERSION` file: rejected because it duplicates Git tag
    state and invites drift.
  - Derive versions from GitHub Releases instead of tags: rejected because tags
    are cheaper to inspect locally and already drive the current release helper.

## Decision 4: Use semver tags on the commit as the idempotency source of truth

- **Decision**: Serialize release runs for `main` and treat a reachable semver
  tag that already points at the pushed commit as the idempotency source of
  truth so reruns or overlapping pushes cannot create duplicate releases.
- **Rationale**: Automatic release triggers increase the chance of race
  conditions and reruns, so the workflow must protect release state with a
  single canonical signal. Tags are already the release source of truth, so a
  tag pointing at the commit is safer than inferring state from the GitHub
  release record alone.
- **Alternatives considered**:
  - Allow concurrent release jobs and trust tag creation to fail one of them:
    rejected because it creates noisy partial states and harder-to-debug logs.
  - Use GitHub release records rather than tags to detect reruns: rejected
    because release metadata can drift from the canonical Git tag state.
  - Always fail reruns on already released commits: rejected because a clear
    no-op or skip result is safer and easier for maintainers to interpret.

## Decision 5: Treat tag-without-release drift as a manual repair case

- **Decision**: If a semver tag already points at the commit but the
  corresponding GitHub release record is missing or incomplete, the workflow
  skips automatic publication and asks the maintainer to repair the release
  state manually.
- **Rationale**: Automatic recreation could mask partial-release problems or
  attach artifacts to the wrong history. A visible skip keeps the canonical tag
  state intact and avoids publishing duplicate versions.
- **Alternatives considered**:
  - Recreate the GitHub release automatically from the existing tag: rejected
    because it blurs the line between idempotent reruns and manual recovery.
  - Delete and recreate the tag automatically: rejected because it is
    destructive and risky in CI.

## Decision 6: Publish only after validation and all builds succeed

- **Decision**: The workflow must compute the next minor version, run tests,
  build the full artifact matrix, verify at least one built binary, and only
  then create the tag and GitHub release.
- **Rationale**: This prevents partial releases where a tag exists but artifacts
  or verification failed, which would make automatic publication less
  trustworthy than the current manual process.
- **Alternatives considered**:
  - Create the tag before builds start: rejected because failed builds would
    leave behind misleading release state.
  - Publish artifacts incrementally by platform: rejected because automatic
    releases should remain all-or-nothing.

## Decision 7: Fan out package publishing from the successful GitHub release

- **Decision**: Keep GitHub release build/publish in the main release workflow,
  and let downstream package-manager workflows such as AUR trigger from the
  published-release event.
- **Rationale**: Package managers should consume the verified release artifacts
  rather than duplicate upstream build logic, and a release-published fan-out
  point cleanly supports future package managers.
- **Alternatives considered**:
  - Publish AUR directly inside the main release workflow: rejected because it
    couples downstream packaging secrets and failure modes to the core build.
  - Trigger package-manager publishing from raw `push` events: rejected because
    downstream publishers should only run after a successful GitHub release.

## Decision 8: Keep local preview and documentation aligned with the automated flow

- **Decision**: Maintain a repository-local preview helper and update README plus
  quickstart guidance so contributors can predict the next minor version,
  observe the automatic workflow, and verify published artifacts.
- **Rationale**: Automatic releases are only understandable if contributors can
  reproduce the version calculation locally and recognize the CI behavior from
  the documentation.
- **Alternatives considered**:
  - Document only the workflow file: rejected because contributors still need a
    fast local preview path.
  - Remove the preview helper and rely on reading tags manually: rejected
    because that makes CI behavior harder to validate before pushing.
