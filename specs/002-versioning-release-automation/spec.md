# Feature Specification: Version reporting and release automation

**Feature Branch**: `002-versioning-release-automation`  
**Created**: 2026-03-30  
**Status**: Draft  
**Input**: User description: "I want versioning to be handled properly for this project. Create `ds --version` (or a `version` command) so I can tell which version I'm on, and add automated version bumping through GitHub CI/CD."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Inspect the installed CLI version (Priority: P1)

As a developer or operator, I want a fast way to inspect the installed `ds`
version so I can confirm which build I am running before debugging, onboarding,
or reporting issues.

**Why this priority**: Version visibility is the foundation for all other
release work. If users cannot identify the binary they have installed, releases
and bug reports remain ambiguous.

**Independent Test**: Build `ds` once with injected version metadata and once as
a local development build, then verify `ds --version` and `ds version` return
the expected values and exit successfully without touching any `.env` files.

**Acceptance Scenarios**:

1. **Given** a release binary built with version metadata, **When** the user runs
   `ds --version`, **Then** the CLI prints a concise version string and exits
   with code `0`.
2. **Given** a release binary built with version metadata, **When** the user runs
   `ds version`, **Then** the CLI prints detailed build metadata including the
   semantic version, commit, and build timestamp.
3. **Given** a local development build with no injected release metadata,
   **When** the user runs either version path, **Then** the CLI reports safe
   fallback metadata such as `dev` rather than failing.

---

### User Story 2 - Publish semver releases from GitHub Actions (Priority: P2)

As a maintainer, I want to trigger a release workflow from GitHub Actions with a
major, minor, or patch bump so I can cut consistent releases without manually
editing source files or tagging from my laptop.

**Why this priority**: Manual tagging and source edits are error-prone. A single
release path makes versioning repeatable and keeps release state aligned with
the repository.

**Independent Test**: Trigger the release workflow in a test repository with
`patch`, `minor`, and `major` inputs, then confirm the workflow calculates the
next semantic version from existing tags, runs validation, and publishes a
GitHub release only when all checks pass.

**Acceptance Scenarios**:

1. **Given** the latest reachable semver tag is `v0.4.2`, **When** a maintainer
   triggers the workflow with a `minor` bump, **Then** the workflow computes
   `v0.5.0`, validates the project, and publishes that version.
2. **Given** the repository has no prior semver tags, **When** a maintainer
   triggers the workflow with a bump input, **Then** the workflow derives the
   first release from a `v0.0.0` baseline and publishes the correct next
   version.
3. **Given** tests or packaging fail, **When** the workflow runs, **Then** no
   new tag or release is published and the maintainer sees actionable failure
   logs.

---

### User Story 3 - Build repeatable versioned binaries locally and in CI (Priority: P3)

As a contributor, I want a documented and repeatable build path for versioned
artifacts so local validation, CI builds, and GitHub releases all report
matching metadata.

**Why this priority**: Release automation is only trustworthy if local builds,
CI builds, and published artifacts all use the same metadata contract.

**Independent Test**: Follow the documented local build steps to inject
version/commit/build-time values, run `ds version`, and compare the reported
metadata with the GitHub release artifact contract.

**Acceptance Scenarios**:

1. **Given** explicit build metadata values, **When** a contributor builds `ds`
   locally with the documented `ldflags`, **Then** the resulting binary reports
   those exact values through the version command.
2. **Given** the release workflow publishes artifacts for the supported target
   matrix, **When** those binaries are executed with `--version`, **Then** they
   all report the same semantic version for that release.
3. **Given** the README and quickstart are updated, **When** a contributor
   follows them, **Then** they can inspect versions and understand the release
   flow without reading implementation code.

### Edge Cases

- What happens when the repository contains non-semver tags that should not
  influence the next release version?
- What happens when the requested release version already exists because a prior
  workflow run partially completed?
- How does the CLI report version information when the binary was built outside
  Git metadata or without explicit `ldflags`?
- What happens if the workflow is triggered from a non-default branch or against
  a stale checkout?
- How are partially built artifacts handled if one target fails after others
  succeed?

## User Experience Consistency _(mandatory)_

- **UX-001**: `ds --version` and `ds version` MUST use stable, documented output
  formats so users can quickly identify the installed build.
- **UX-002**: Release workflow failures MUST explain the blocking condition, its
  impact on publishing, and the next maintainer action.
- **UX-003**: Version and release output MUST never print secrets or raw tokens;
  logs may show release metadata only.
- **UX-004**: Tag names, GitHub release titles, and artifact file names MUST
  match the published semantic version exactly so users can trust what they
  download.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST support a root `ds --version` path for fast version
  inspection.
- **FR-002**: The system MUST support a `ds version` command for detailed build
  metadata inspection.
- **FR-003**: Version inspection MUST report the semantic version string and MUST
  provide commit and build-time metadata in the detailed command output.
- **FR-004**: Local development builds without injected metadata MUST still
  return deterministic fallback values and exit successfully.
- **FR-005**: Build metadata MUST be injected at build time rather than requiring
  maintainers to edit Go source files for each release.
- **FR-006**: The system MUST provide a GitHub Actions `workflow_dispatch`
  release workflow with a required `major`, `minor`, or `patch` bump input and
  optional release notes input.
- **FR-007**: The release workflow MUST derive the next semantic version from the
  latest reachable semver tag, ignoring unrelated tags and using `v0.0.0` as the
  baseline when no semver tag exists.
- **FR-008**: The release workflow MUST validate the repository from the default
  branch and MUST fail before publishing if tests, version calculation, or
  packaging fail.
- **FR-009**: The release workflow MUST publish a Git tag, a GitHub release, and
  versioned `ds` binaries for the supported cross-platform target matrix.
- **FR-010**: Published artifact names and release titles MUST include the exact
  semantic version they contain.
- **FR-011**: README and quickstart documentation MUST explain version
  inspection, local metadata injection, and the maintainer release workflow.
- **FR-012**: Existing test automation (`go test ./...`) MUST remain the
  validation gate for release publication.

### Key Entities _(include if feature involves data)_

- **Version Metadata**: The build-identifying values embedded into the binary,
  including semantic version, commit SHA, build timestamp, and runtime platform.
- **Release Request**: The maintainer-provided release intent captured by the
  workflow input, including bump level, target ref, and optional notes.
- **Release Artifact**: A packaged `ds` binary produced for a specific operating
  system and architecture and associated with a single semantic version.
- **Release Publication**: The published tag and GitHub release record that ties
  together the computed version, commit, notes, and attached artifacts.

## Performance Requirements _(mandatory)_

- **PR-001**: `ds --version` and `ds version` MUST complete within 100 ms on a
  normal local machine because they only read embedded metadata.
- **PR-002**: The versioning feature MUST not add any new runtime dependency to
  the shipped `ds` binary beyond the existing Go module dependencies.
- **PR-003**: The manual release workflow MUST complete validation, cross-build,
  and publication for the supported target matrix within 15 minutes under normal
  GitHub-hosted runner availability.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Maintainers can identify the installed `ds` build from the command
  line in one command with no repository access.
- **SC-002**: A maintainer can publish a patch, minor, or major release without
  manually editing Go source files or creating Git tags locally.
- **SC-003**: 100% of published release artifacts report the exact release
  version through the CLI version paths.
- **SC-004**: Release attempts that fail validation or packaging produce no
  orphaned published tag or GitHub release.
