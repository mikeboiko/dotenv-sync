# Data Model: Automatic patch release automation

## Release Trigger

- **Purpose**: Represents the push event that may start an automatic release.
- **Fields**:
  - `event_name`: workflow event type, expected to be `push`
  - `branch`: pushed branch name, expected to be `main`
  - `commit`: full Git SHA for the pushed commit
  - `actor`: GitHub user or app that pushed the commit
  - `already_released`: boolean indicating whether a reachable semver tag already
    points at the commit
- **Validation rules**:
  - `event_name` must be `push`
  - `branch` must equal the repository default branch before publication begins
  - `commit` must be non-empty
  - `already_released` is derived from Git tag state, not GitHub release
    metadata

## Release Publication

- **Purpose**: Represents the computed release outcome before and after
  publication.
- **Fields**:
  - `previous_version`: latest reachable semver tag or `v0.0.0`
  - `next_version`: computed patch version to publish
  - `commit`: commit SHA included in the release build
  - `status`: enum such as `planned`, `validated`, `built`, `published`,
    `skipped`, or `failed`
  - `skip_reason`: optional explanation when publication is intentionally skipped
  - `release_url`: GitHub release URL after publication
- **Validation rules**:
  - `next_version` must be strictly greater than `previous_version` when
    `status` reaches `published`
  - a single `commit` may have at most one `published` release
  - `skip_reason` is required when `status` is `skipped`
  - if a semver tag already points at `commit`, the workflow treats the
    publication as `skipped` even when a GitHub release record needs manual
    repair

## Release Artifact

- **Purpose**: Represents one packaged binary attached to a release.
- **Fields**:
  - `version`: semantic version carried by the binary
  - `os`: target operating system
  - `arch`: target architecture
  - `archive_name`: published asset file name
  - `checksum`: integrity value for the archive
  - `published`: boolean indicating whether the asset is attached to the release
- **Validation rules**:
  - `archive_name` must include the exact release version
  - all artifacts for one publication must share the same `version`

## Version Metadata

- **Purpose**: Represents the build-identifying values embedded into a `ds`
  binary and checked during release verification.
- **Fields**:
  - `version`: semantic version string with a leading `v`
  - `commit`: Git SHA embedded into the binary
  - `build_time`: UTC build timestamp string
  - `platform`: runtime platform string such as `linux/amd64`
- **Validation rules**:
  - `version` must equal the release tag for published artifacts
  - `build_time` must be machine-readable when injected

## Relationships

- One **Release Trigger** produces zero or one **Release Publication**.
- One **Release Publication** may contain many **Release Artifacts**.
- One published **Release Publication** must correspond to exactly one `main`
  commit.
- Every published **Release Artifact** must embed the **Version Metadata** of its
  parent **Release Publication**.

## State Transitions

- **Release Publication** success path: `planned -> validated -> built -> published`
- **Release Publication** skip path: `planned -> skipped`
- **Release Publication** failure path: `planned|validated|built -> failed`
- **Release Artifact**: `packaged -> attached`
