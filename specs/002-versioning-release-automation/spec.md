# Feature Specification: Automatic patch release automation

**Feature Branch**: `002-versioning-release-automation`  
**Created**: 2026-03-30  
**Status**: Draft  
**Input**: User description: "Right now, I have a manual release CI/CD process. I want you to automate this, so a patch version is bumped, every time I push to main."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Publish the next patch release on every push to `main` (Priority: P1)

As a maintainer, I want each push to `main` to publish the next patch release
automatically so release versions stay aligned with merged code without manual
workflow dispatches or local tagging.

**Why this priority**: This is the direct user request and the smallest slice
that turns the current manual release flow into automatic CI/CD.

**Independent Test**: Push a new commit to `main` in a test repository where the
latest semver tag is `v0.4.2`, then verify the workflow computes `v0.4.3`, runs
validation, and publishes tagged artifacts only after the build matrix succeeds.

**Acceptance Scenarios**:

1. **Given** the latest reachable semver tag is `v0.4.2`, **When** a new commit
   is pushed to `main`, **Then** the workflow computes `v0.4.3`, validates the
   repository, and publishes that release.
2. **Given** the repository has no prior semver tags, **When** the first commit
   is pushed to `main` after automation is enabled, **Then** the workflow uses
   `v0.0.0` as the baseline and publishes `v0.0.1`.
3. **Given** two different commits are pushed to `main` in sequence, **When**
   both workflows complete successfully, **Then** each commit receives a unique,
   monotonically increasing patch release with matching tag, release title, and
   artifact names.

---

### User Story 2 - Keep automatic releases safe, idempotent, and actionable (Priority: P2)

As a maintainer, I want reruns and failure paths to stop cleanly without creating
duplicate or partial releases so the automated pipeline remains trustworthy.

**Why this priority**: Automatic release creation is only safe if reruns,
validation failures, and non-semver repository noise do not corrupt release
state.

**Independent Test**: Rerun the workflow for an already released `main` commit,
then simulate failing tests or packaging and confirm that no new tag or GitHub
release is created while logs explain the blocking condition.

**Acceptance Scenarios**:

1. **Given** a `main` commit that already has a reachable semver tag pointing at
   it, **When** the workflow is rerun for that same commit, **Then** the
   workflow exits without publishing another version, reports that the commit is
   already released by tag, and requires manual maintainer repair if the GitHub
   release record is missing.
2. **Given** `go test ./...` or artifact packaging fails, **When** the workflow
   runs for a new `main` commit, **Then** no new tag or release is published and
   the failure logs explain what blocked publication and what to inspect next.
3. **Given** unrelated Git tags exist in the repository, **When** the workflow
   calculates the next release, **Then** those tags are ignored and only the
   latest reachable semantic version influences the patch bump.

---

### User Story 3 - Preview and verify automatic releases locally and in CI (Priority: P3)

As a contributor, I want a local preview and verification path that matches CI so
I can predict the next automatic patch release and confirm published binaries
self-report the expected version.

**Why this priority**: Automatic releases are easier to trust when contributors
can reproduce the version calculation locally and compare it to CI output.

**Independent Test**: Run the local preview helper against a repository with
existing tags, build a binary with the predicted version metadata, and compare
its `ds --version` output with the release artifact contract used in CI.

**Acceptance Scenarios**:

1. **Given** the latest semver tag is `v0.4.2`, **When** a contributor runs the
   local preview helper, **Then** it reports `v0.4.3` without requiring a manual
   bump argument.
2. **Given** the release workflow publishes artifacts for the supported target
   matrix, **When** one of those binaries is executed with `--version`, **Then**
   it reports the exact semantic version carried by the tag and release title.
3. **Given** the README and quickstart are updated, **When** a contributor
   follows them, **Then** they can understand the push-to-`main` release flow,
   how to monitor it, and how to verify a published artifact without reading the
   implementation.

### Edge Cases

- What happens when `main` receives multiple pushes before an earlier release run
  completes?
- What happens when the workflow is rerun for a commit that already has a
  reachable semver tag?
- What happens when the repository contains non-semver tags that should not
  influence the next patch version?
- What happens when the repository has no prior semver tags?
- How are partially built artifacts handled if one target fails after others
  succeed?
- What happens when a push occurs on a branch other than `main`?

## User Experience Consistency _(mandatory)_

- **UX-001**: Workflow logs MUST clearly distinguish validation, skip, failure,
  and publish outcomes for automatic releases.
- **UX-002**: Failure and skip messages MUST explain the blocking condition, its
  impact on publication, and the next maintainer action.
- **UX-003**: Release logs, preview output, and verification steps MUST never
  print secrets or raw tokens; they may show refs, versions, and artifact names
  only.
- **UX-004**: Tag names, GitHub release titles, and artifact file names MUST
  match the published semantic version exactly.
- **UX-005**: Local preview and CI automation MUST describe the behavior as a
  patch-only release flow triggered by pushes to `main`.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST trigger release automation automatically on pushes
  to the repository default branch (`main`).
- **FR-002**: The system MUST compute the next release version as the next patch
  increment from the latest reachable semantic version tag, using `v0.0.0` as
  the baseline when no semver tag exists.
- **FR-003**: The system MUST not require manual bump selection, manual source
  edits, or local Git tag creation to publish a patch release.
- **FR-004**: The workflow MUST validate the pushed `main` commit and run
  `go test ./...` before creating any release tag or GitHub release.
- **FR-005**: The workflow MUST publish a Git tag, a GitHub release, and
  versioned `ds` binaries for the supported target matrix only after all builds
  succeed.
- **FR-006**: The system MUST publish at most one semantic version per `main`
  commit and MUST treat a reachable semver tag on that commit as the idempotency
  source of truth during reruns or repeated workflow execution.
- **FR-007**: The system MUST ignore non-semver tags when computing the next
  patch release.
- **FR-008**: The system MUST serialize automatic release runs for `main` so
  overlapping pushes cannot publish conflicting versions.
- **FR-009**: Published artifact names and GitHub release titles MUST include the
  exact semantic version they contain.
- **FR-010**: The workflow MUST directly verify at least one built Linux
  reference artifact with the existing `ds --version` or `ds version` path
  before publication, while automated tests enforce version-parity rules for the
  remaining release artifacts.
- **FR-011**: Local preview tooling and documentation MUST explain how to
  predict the next automatic patch release and how to monitor the release run.
- **FR-012**: Existing test automation (`go test ./...`) MUST remain the
  validation gate for release publication.

### Key Entities _(include if feature involves data)_

- **Release Trigger**: The push event that may start automatic publication,
  including branch, commit SHA, actor, and whether that commit is already
  released.
- **Release Publication**: The computed release outcome before and after
  publication, including the previous version, next version, commit SHA, status,
  and optional skip or failure reason.
- **Release Artifact**: A packaged `ds` binary produced for a specific operating
  system and architecture and associated with a single semantic version.
- **Version Metadata**: The build-identifying values embedded into a binary and
  checked during CI verification.

## Performance Requirements _(mandatory)_

- **PR-001**: Automatic release publication for the supported target matrix MUST
  complete within 15 minutes under normal GitHub-hosted runner availability.
- **PR-002**: Local patch preview MUST complete within 1 second on a normal
  development machine because it only inspects repository metadata.
- **PR-003**: The feature MUST not add any new runtime dependency to the shipped
  `ds` binary beyond the existing Go module dependencies.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Maintainers can publish a new patch release by pushing to `main`
  without manually dispatching a workflow or creating a Git tag locally.
- **SC-002**: Every successful `main` push that reaches the release workflow can
  produce exactly one new patch release with matching tag, release title, and
  artifact names.
- **SC-003**: 100% of published release artifacts are built from the same release
  version input, with the Linux reference artifact verified directly in CI and
  cross-platform version parity enforced by automated tests.
- **SC-004**: Workflow reruns or failed validations publish no duplicate or
  partial releases.
- **SC-005**: Contributors can predict the next automatic patch version locally
  and follow the documented monitoring and verification flow without reading the
  implementation.
