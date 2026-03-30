# Data Model: Version reporting and release automation

## Version Metadata

- **Purpose**: Represents the build-identifying values embedded into a `ds`
  binary and exposed through the CLI.
- **Fields**:
  - `version`: semantic version string with `v` prefix for releases, or `dev`
    for local builds
  - `commit`: Git SHA or fallback placeholder such as `none`
  - `build_time`: UTC timestamp string or fallback placeholder such as `unknown`
  - `platform`: runtime platform string such as `linux/amd64`
  - `source`: enum describing whether metadata came from injected release values
    or local defaults
- **Validation rules**:
  - `version` must be non-empty
  - release versions must follow semantic versioning with a leading `v`
  - `build_time` must be machine-readable when injected

## Release Request

- **Purpose**: Captures the maintainer input to the manual release workflow.
- **Fields**:
  - `bump`: enum with `major`, `minor`, or `patch`
  - `target_ref`: Git ref the workflow is releasing from
  - `notes`: optional release notes override
  - `requested_by`: GitHub actor that triggered the workflow
- **Validation rules**:
  - `bump` must be one of the supported semver levels
  - `target_ref` must resolve to the default branch head before publication

## Release Publication

- **Purpose**: Represents the computed release outcome before and after
  publishing.
- **Fields**:
  - `previous_version`: latest reachable semver tag or `v0.0.0`
  - `next_version`: computed semantic version to publish
  - `tag_name`: Git tag to create, identical to `next_version`
  - `commit`: commit SHA included in the release build
  - `status`: enum such as `planned`, `validated`, `built`, `published`, or
    `failed`
  - `release_url`: GitHub release URL after publication
- **Validation rules**:
  - `next_version` must be strictly greater than `previous_version`
  - `tag_name` must not already exist before publication starts

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

## Relationships

- One **Release Request** produces zero or one **Release Publication**.
- One **Release Publication** may contain many **Release Artifacts**.
- One built `ds` binary exposes one **Version Metadata** record at runtime.
- Every **Release Artifact** must embed the **Version Metadata** associated with
  its parent **Release Publication**.

## State Transitions

- **Release Publication**: `planned -> validated -> built -> published`
- **Release Publication** failure path: `planned|validated|built -> failed`
- **Release Artifact**: `packaged -> attached`
