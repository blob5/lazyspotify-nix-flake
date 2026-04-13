#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

source_dir=""
output=""
daemon_version="${DAEMON_VERSION:-}"

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
    --daemon-version)
      daemon_version="$2"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ -z "${source_dir}" ]]; then
  echo "--source-dir is required" >&2
  exit 1
fi

if [[ -z "${output}" ]]; then
  echo "--output is required" >&2
  exit 1
fi

if [[ -z "${daemon_version}" ]]; then
  daemon_version="$(default_daemon_version)"
fi

mkdir -p "$(dirname "${output}")"

(
  cd "${source_dir}"
  go build \
    -trimpath \
    -buildvcs=false \
    -ldflags "-X github.com/devgianlu/go-librespot.version=${daemon_version}" \
    -o "${output}" \
    ./cmd/daemon
)
