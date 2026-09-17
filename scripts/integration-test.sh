#!/usr/bin/env bash
set -euo pipefail

binary=${1:-./bin/golangci-lint-passthrough}
fixture=./integration/testdata

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

run_fixture() {
  local config=$1
  local output=$2
  set +e
  "$binary" run --config "$config" "$fixture" >"$output" 2>&1
  local status=$?
  set -e
  if [[ $status -ne 1 ]]; then
    cat "$output"
    echo "expected lint findings (exit 1), got exit $status" >&2
    exit 1
  fi
}

run_fixture integration/default.yml "$tmp_dir/default.out"
grep -q 'function ordinary is a passthrough to target' "$tmp_dir/default.out"
if grep -q 'function withExtra' "$tmp_dir/default.out"; then
  cat "$tmp_dir/default.out"
  echo 'default configuration unexpectedly accepted an extra argument' >&2
  exit 1
fi
if grep -q 'function suppressed' "$tmp_dir/default.out"; then
  cat "$tmp_dir/default.out"
  echo '//nolint:passthrough did not suppress the finding' >&2
  exit 1
fi

run_fixture integration/nondefault.yml "$tmp_dir/nondefault.out"
grep -q 'function ordinary is a passthrough to target' "$tmp_dir/nondefault.out"
grep -q 'function withExtra is a passthrough to target' "$tmp_dir/nondefault.out"
if grep -q 'function suppressed' "$tmp_dir/nondefault.out"; then
  cat "$tmp_dir/nondefault.out"
  echo '//nolint:passthrough did not suppress the finding' >&2
  exit 1
fi

echo 'integration fixtures passed'
