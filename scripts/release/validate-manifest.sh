#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

require_daemon_tag=0
if [[ "${1:-}" == "--require-daemon-tag" ]]; then
  require_daemon_tag=1
fi

version="$(release_version)"
repo="$(daemon_repo)"
tag="$(daemon_tag)"
commit="$(daemon_commit)"
bundle="$(bundle_version)"

if [[ -z "${version}" ]]; then
  echo "lazyspotify_version is required" >&2
  exit 1
fi

if [[ -z "${repo}" ]]; then
  echo "daemon_repo is required" >&2
  exit 1
fi

if [[ -z "${commit}" ]]; then
  echo "daemon_commit is required" >&2
  exit 1
fi

if [[ ! "${commit}" =~ ^[0-9a-f]{40}$ ]]; then
  echo "daemon_commit must be a full 40-character git commit" >&2
  exit 1
fi

if [[ -z "${bundle}" ]]; then
  echo "bundle_version is required" >&2
  exit 1
fi

if (( require_daemon_tag )) && [[ -z "${tag}" ]]; then
  echo "daemon_tag must be set before publishing a release" >&2
  exit 1
fi
