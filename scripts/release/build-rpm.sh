#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

source_bundle=""
output_dir="$(default_dist_dir)"
source_commit="$(current_commit)"
build_date="$(build_date_utc)"
daemon_version="$(default_daemon_version)"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --source-bundle)
      source_bundle="$2"
      shift 2
      ;;
    --output-dir)
      output_dir="$2"
      shift 2
      ;;
    --source-commit)
      source_commit="$2"
      shift 2
      ;;
    --build-date)
      build_date="$2"
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

if [[ -z "${source_bundle}" ]]; then
  echo "--source-bundle is required" >&2
  exit 1
fi

require_command rpmbuild

mkdir -p "${output_dir}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

topdir="${tmpdir}/rpmbuild"
mkdir -p "${topdir}/BUILD" "${topdir}/RPMS" "${topdir}/SOURCES" "${topdir}/SPECS" "${topdir}/SRPMS"
cp "${source_bundle}" "${topdir}/SOURCES/$(source_bundle_filename)"
cp "${REPO_ROOT}/packaging/rpm/lazyspotify.spec" "${topdir}/SPECS/lazyspotify.spec"

rpmbuild -ba \
  --define "_topdir ${topdir}" \
  --define "lazyspotify_version $(release_version)" \
  --define "lazyspotify_release 1" \
  --define "source_commit ${source_commit}" \
  --define "build_date ${build_date}" \
  --define "daemon_version ${daemon_version}" \
  "${topdir}/SPECS/lazyspotify.spec"

find "${topdir}" \
  \( -name '*.rpm' -o -name '*.src.rpm' \) \
  -exec cp {} "${output_dir}/" \;
