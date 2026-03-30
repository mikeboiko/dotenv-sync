# Quickstart: Version reporting and release automation

## Check the current local build

Build the CLI normally:

```bash
go build -o ./bin/ds ./cmd/ds
./bin/ds --version
./bin/ds version
```

Expected development output:

```text
ds dev
Version: dev
Commit: none
Built: unknown
Platform: linux/amd64
```

## Build a local binary with explicit version metadata

```bash
VERSION=v0.1.0
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -o ./bin/ds \
  -ldflags "-X dotenv-sync/pkg/dotenvsync.Version=$VERSION -X dotenv-sync/pkg/dotenvsync.Commit=$COMMIT -X dotenv-sync/pkg/dotenvsync.BuildTime=$BUILD_TIME" \
  ./cmd/ds

./bin/ds --version
./bin/ds version
```

## Install the current checkout into `~/.local/bin/ds`

```bash
git switch main
git pull --ff-only
./scripts/install-local.sh
ds version
```

If `main` is exactly on the latest release tag, the reported version matches the
latest GitHub release. If `main` is ahead of the latest release, the script uses
Git-derived metadata such as `v0.4.0-3-gabc1234` so the local build remains
distinguishable from the last published release.

## Preview the next release version locally

```bash
go run ./scripts/nextversion --bump patch
```

If no prior semver tag exists, the patch baseline is `v0.0.1`.

## Trigger a GitHub release

From the GitHub UI, run the manual release workflow and choose `major`,
`minor`, or `patch`.

The workflow only publishes when it runs from the latest commit on the default
branch (typically `main`).

Or use the GitHub CLI:

```bash
gh workflow run release.yml -f bump=patch -f release_notes="Patch release"
gh run watch
```

## Verify a published release

After the workflow completes:

1. Open the GitHub release for the new tag.
2. Download one of the published `ds_<version>_<os>_<arch>` artifacts.
3. Run `ds --version` or `ds version` and confirm the reported version matches
   the tag name exactly.
