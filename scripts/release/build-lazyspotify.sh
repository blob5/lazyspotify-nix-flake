#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

source_dir="${REPO_ROOT}"
output=""
daemon_path=""
version="$(release_version)"
commit=""
build_date=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --source-dir)
      source_dir="$2"
      shift 2
      ;;
    --output)
      output="$2"
      shift 2
      ;;
    --daemon-path)
      daemon_path="$2"
      shift 2
      ;;
    --version)
      version="$2"
      shift 2
      ;;
    --commit)
      commit="$2"
      shift 2
      ;;
    --build-date)
      build_date="$2"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ -z "${output}" ]]; then
  echo "--output is required" >&2
  exit 1
fi

load_bundle_metadata "${source_dir}"

if [[ -z "${commit}" ]]; then
  commit="${LAZYSPOTIFY_COMMIT:-$(current_commit)}"
fi

if [[ -z "${build_date}" ]]; then
  build_date="$(build_date_utc)"
fi

mkdir -p "$(dirname "${output}")"

(
  cd "${source_dir}"
  GOFLAGS="${GOFLAGS:-}" \
  CGO_ENABLED="${CGO_ENABLED:-0}" \
  go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "$(lazyspotify_ldflags "${version#v}" "${commit}" "${build_date}" "${daemon_path}")" \
    -o "${output}" \
    ./cmd/lazyspotify
)
