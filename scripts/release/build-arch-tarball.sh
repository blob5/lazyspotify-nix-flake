#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

source_dir="${REPO_ROOT}"
daemon_source_dir=""
output_dir="$(default_dist_dir)"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --source-dir)
      source_dir="$2"
      shift 2
      ;;
    --daemon-source-dir)
      daemon_source_dir="$2"
      shift 2
      ;;
    --output-dir)
      output_dir="$2"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ -z "${daemon_source_dir}" ]]; then
  daemon_source_dir="${source_dir}/third_party/go-librespot"
fi

mkdir -p "${output_dir}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

artifact_dir="${tmpdir}/$(artifact_prefix)-arch-amd64"
GOOS=linux GOARCH=amd64 \
  "$(dirname "${BASH_SOURCE[0]}")/build-lazyspotify.sh" \
    --source-dir "${source_dir}" \
    --output "${artifact_dir}/lazyspotify" \
    --daemon-path "/usr/lib/lazyspotify/lazyspotify-librespot"

GOOS=linux GOARCH=amd64 \
  "$(dirname "${BASH_SOURCE[0]}")/build-librespot-daemon.sh" \
    --source-dir "${daemon_source_dir}" \
    --output "${artifact_dir}/lazyspotify-librespot"

cp "${source_dir}/LICENSE" "${artifact_dir}/LICENSE"
cp "${daemon_source_dir}/LICENSE" "${artifact_dir}/LICENSE.go-librespot"
cp "${source_dir}/README.md" "${artifact_dir}/README.md"

tar -C "${tmpdir}" -czf "${output_dir}/$(artifact_prefix)-arch-amd64.tar.gz" "$(artifact_prefix)-arch-amd64"
printf '%s\n' "${output_dir}/$(artifact_prefix)-arch-amd64.tar.gz"
