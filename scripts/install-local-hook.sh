#!/usr/bin/env bash

set -euo pipefail

script_dir="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

default_branch="$(git -C "$repo_root" symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null || true)"
default_branch="${default_branch#origin/}"
if [[ -z "$default_branch" ]]; then
    default_branch="main"
fi

current_branch="$(git -C "$repo_root" branch --show-current 2>/dev/null || true)"
if [[ -z "$current_branch" || "$current_branch" != "$default_branch" ]]; then
    exit 0
fi

exec "${script_dir}/install-local.sh" --quiet
