# Workflow Contract: Automatic patch release automation

## Workflow

- **Path**: `.github/workflows/release.yml`
- **Trigger**: `push` to the repository default branch (`main`)

## Inputs

- No manual release inputs are required.
- Version calculation is always patch-only.

## Preconditions

- The workflow must run only for pushes to the repository default branch head.
- The workflow must have permission to create Git tags and GitHub releases.
- The workflow must serialize release runs for `main` to avoid conflicting
  version publication.
- The repository must pass `go test ./...` before publication begins.

## Behavior

1. Determine the latest reachable semantic version tag, ignoring unrelated tags.
2. Compute the next patch version using `v0.0.0` as the baseline when no semver
   tag exists.
3. If a reachable semver tag already points at the pushed commit, exit without
   creating another tag or release, report a clear skip reason, and require
   manual maintainer repair if the GitHub release record is missing.
4. Run validation and tests before creating any tag or GitHub release.
5. Build the supported `ds` target matrix with embedded version metadata.
6. Package release artifacts using deterministic versioned names.
7. Verify the Linux reference artifact directly with `ds --version` before
   publication, and rely on automated tests to enforce version parity across the
   remaining artifact matrix.
8. Create the Git tag and publish the GitHub release only after all builds
   succeed.

## Artifact Naming

- POSIX archives should follow:

```text
ds_<version>_<os>_<arch>.tar.gz
```

- Windows archives should follow:

```text
ds_<version>_windows_<arch>.zip
```

- Every asset name must include the exact semantic version it contains.

## Failure and Skip Rules

- If version calculation, validation, or any build fails, the workflow must stop
  before creating the release tag.
- If a rerun targets a commit that already has a reachable semver tag, the
  workflow must skip publication, explain why, and avoid recreating release
  state automatically.
- If overlapping `main` pushes occur, release runs must remain serialized so
  version publication cannot race.
- Logs may show version numbers, refs, and artifact names, but must not expose
  tokens or secrets.
