# Workflow Contract: Manual release automation

## Workflow

- **Path**: `.github/workflows/release.yml`
- **Trigger**: `workflow_dispatch`

## Inputs

| Input           | Required | Allowed values            | Purpose                                                   |
| --------------- | -------- | ------------------------- | --------------------------------------------------------- |
| `bump`          | Yes      | `major`, `minor`, `patch` | Select the next semantic-version increment                |
| `release_notes` | No       | freeform text             | Optional maintainer-supplied notes for the GitHub release |

## Preconditions

- The workflow must release from the repository's default branch head.
- The workflow must have permission to create Git tags and GitHub releases.
- The repository must pass `go test ./...` before publication begins.

## Behavior

1. Determine the latest reachable semantic version tag, ignoring unrelated tags.
2. Compute the next semantic version from the requested bump using `v0.0.0` as
   the baseline when no semver tag exists.
3. Run validation and tests before creating any tag or GitHub release.
4. Build the supported `ds` target matrix with embedded version metadata.
5. Package release artifacts using deterministic versioned names.
6. Create the Git tag and publish the GitHub release only after all builds
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

## Failure Rules

- If version calculation, validation, or any build fails, the workflow must stop
  before creating the release tag.
- If the computed tag already exists, the workflow must fail with an actionable
  explanation instead of overwriting release state.
- Failure logs may show version numbers, refs, and artifact names, but must not
  expose tokens or secrets.
