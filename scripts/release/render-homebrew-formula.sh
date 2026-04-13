#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

source_bundle=""
output=""
build_date="$(build_date_utc)"
commit="$(current_commit)"
daemon_version="$(default_daemon_version)"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --source-bundle)
      source_bundle="$2"
      shift 2
      ;;
    --output)
      output="$2"
      shift 2
      ;;
    --build-date)
      build_date="$2"
      shift 2
      ;;
    --commit)
      commit="$2"
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

if [[ -z "${source_bundle}" || -z "${output}" ]]; then
  echo "--source-bundle and --output are required" >&2
  exit 1
fi

source_sha="$(sha256_file "${source_bundle}")"
source_url="$(release_asset_url "$(source_bundle_filename)")"
mkdir -p "$(dirname "${output}")"

sed \
  -e "s|@VERSION@|$(release_version)|g" \
  -e "s|@SOURCE_URL@|${source_url}|g" \
  -e "s|@SOURCE_SHA256@|${source_sha}|g" \
  -e "s|@COMMIT@|${commit}|g" \
  -e "s|@BUILD_DATE@|${build_date}|g" \
  -e "s|@DAEMON_VERSION@|${daemon_version}|g" \
  "${REPO_ROOT}/packaging/homebrew/lazyspotify.rb.tmpl" > "${output}"
