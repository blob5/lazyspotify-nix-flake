#!/usr/bin/env bash

set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

output_dir="${1:-$(default_dist_dir)}"
output_file="${output_dir}/SHA256SUMS"

mkdir -p "${output_dir}"
: > "${output_file}"

while IFS= read -r path; do
  printf '%s  %s\n' "$(sha256_file "${path}")" "$(basename "${path}")" >> "${output_file}"
done < <(find "${output_dir}" -maxdepth 1 -type f ! -name 'SHA256SUMS' | LC_ALL=C sort)

printf '%s\n' "${output_file}"
