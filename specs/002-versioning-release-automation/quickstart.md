# Quickstart: Automatic patch release automation

## Preview the next automatic patch release locally

Run the local preview helper from the repository root:

```bash
go run ./scripts/nextversion
```

Expected examples:

```text
v0.0.1   # when no prior semver tags exist
v0.4.3   # when the latest reachable semver tag is v0.4.2
```

## Build a local binary with the predicted release metadata

```bash
VERSION=$(go run ./scripts/nextversion)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -o ./bin/ds \
  -ldflags "-X dotenv-sync/pkg/dotenvsync.Version=$VERSION -X dotenv-sync/pkg/dotenvsync.Commit=$COMMIT -X dotenv-sync/pkg/dotenvsync.BuildTime=$BUILD_TIME" \
  ./cmd/ds

./bin/ds --version
./bin/ds version
```

This mirrors the metadata contract the release workflow verifies before
publication.

## Push to `main` to trigger an automatic patch release

```bash
git switch main
git pull --ff-only
git push origin main
```

Monitor the workflow from the GitHub UI, or with the GitHub CLI:

```bash
gh run list --workflow release.yml --limit 1
gh run watch <run-id>
```

Under normal GitHub-hosted runner availability, the release run should finish
within 15 minutes.

The workflow should calculate the next patch version, run `go test ./...`, build
the release artifacts, verify the Linux reference artifact with `ds --version`,
and rely on automated tests for cross-platform version parity before publishing.

## Verify a published release

After the workflow completes:

1. Open the GitHub release for the new tag.
2. Download one of the published `ds_<version>_<os>_<arch>` artifacts.
3. Run `ds --version` or `ds version` and confirm the reported version matches
   the tag name exactly.

## Confirm rerun safety

If the release workflow is rerun for a commit that already has a reachable
semver tag, it should report that the commit is already released by tag and
avoid creating another tag or GitHub release. If the tag exists but the GitHub
release record is missing, repair that drift manually instead of expecting the
workflow to republish it automatically.
