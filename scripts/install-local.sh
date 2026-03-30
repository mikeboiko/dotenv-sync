#!/usr/bin/env bash

set -euo pipefail

usage() {
    cat <<'EOF'
Usage: ./scripts/install-local.sh [--bin /absolute/path/to/ds] [--quiet]

Build and install ds from the current checkout.

By default this installs to:
  ~/.local/bin/ds

Version metadata is derived from the current Git checkout:
  - exact release tag: v1.2.3
  - commits ahead of a tag: v1.2.3-4-gabc1234
  - no semver tags yet: dev-abc1234

If you want the exact same version as the latest GitHub release, build from the
release tag (or install the GitHub release artifact). Building from main after
new commits will intentionally show a newer local version string.
EOF
}

normalize_version() {
    local value="${1:-}"
    if [[ -z "$value" ]]; then
        printf 'dev\n'
        return
    fi

    if [[ "$value" =~ ^[0-9a-f]{7,}(-dirty)?$ ]]; then
        printf 'dev-%s\n' "$value"
        return
    fi

    printf '%s\n' "$value"
}

script_dir="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"
bin_path="${HOME}/.local/bin/ds"
quiet=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --bin)
            if [[ $# -lt 2 ]]; then
                echo "missing value for --bin" >&2
                exit 1
            fi
            bin_path="$2"
            shift 2
            ;;
        --quiet)
            quiet=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "unknown argument: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

git_version="$(git -C "$repo_root" describe --tags --match 'v[0-9]*.[0-9]*.[0-9]*' --dirty --always 2>/dev/null || true)"
version="$(normalize_version "$git_version")"
commit="$(git -C "$repo_root" rev-parse --short HEAD 2>/dev/null || printf 'none\n')"
build_time="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
latest_release="$(git -C "$repo_root" tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-version:refname | head -n1 || true)"

mkdir -p "$(dirname "$bin_path")"

(
    cd "$repo_root"
    go build -trimpath \
        -ldflags "-X dotenv-sync/pkg/dotenvsync.Version=${version} -X dotenv-sync/pkg/dotenvsync.Commit=${commit} -X dotenv-sync/pkg/dotenvsync.BuildTime=${build_time}" \
        -o "$bin_path" \
        ./cmd/ds
)

if [[ "$quiet" == "true" ]]; then
    exit 0
fi

printf 'Installed ds to %s\n' "$bin_path"
"$bin_path" version

if [[ -n "$latest_release" ]]; then
    if [[ "$version" == "$latest_release" ]]; then
        printf '\nThis build matches the latest release tag: %s\n' "$latest_release"
    else
        printf '\nLatest release tag: %s\n' "$latest_release"
        printf 'This local build was produced from your current checkout, so it may be newer than the latest GitHub release.\n'
    fi
else
    printf '\nNo semantic release tags were found yet; local installs use repository metadata.\n'
fi
