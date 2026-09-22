#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
CHART="${SCRIPT_DIR}/../charts/pega"
RENDERED="${TMPDIR:-/tmp}/pega-serviceaccount-test.$$"
ERRORS="${RENDERED}.err"
trap 'rm -f "${RENDERED}" "${ERRORS}"' EXIT

pass() {
  printf 'PASS: %s\n' "$1"
}

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

contains() {
  grep -Eq -- "$1" "${RENDERED}"
}

expect_failure() {
  local description="$1"
  shift
  if helm template pega "${CHART}" --namespace pega --set 'global.provider=k8s' "$@" \
      >"${RENDERED}" 2>"${ERRORS}"; then
    fail "${description} unexpectedly rendered successfully"
  fi
  pass "${description}"
}

command -v helm >/dev/null 2>&1 || fail "helm is required"

helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' >"${RENDERED}"
if contains '^kind: ServiceAccount$'; then
  fail "ServiceAccount rendered while management is disabled"
fi
pass "ServiceAccount management is disabled by default"

helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'serviceAccount.enabled=true' \
  --set 'serviceAccount.create=true' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=true' \
  --set 'global.actions.execute=deploy' >"${RENDERED}"
contains 'name: pega-serviceaccount' || fail "runtime ServiceAccount is missing"
contains 'name: pega-installer-serviceaccount' || fail "installer ServiceAccount is missing"
contains 'automountServiceAccountToken: false' || fail "token automount is not disabled"
pass "Managed runtime and installer ServiceAccounts"

helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'serviceAccount.enabled=true' \
  --set 'serviceAccount.create=false' \
  --set 'serviceAccount.name=existing-runtime' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=false' \
  --set 'installer.serviceAccount.name=existing-installer' \
  --set 'global.actions.execute=deploy' >"${RENDERED}"
if contains '^kind: ServiceAccount$'; then
  fail "external ServiceAccounts were unexpectedly created"
fi
contains 'serviceAccountName: existing-runtime' || fail "existing runtime account is not assigned"
pass "Pre-existing ServiceAccounts"

helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=false' \
  --set 'installer.serviceAccount.name=existing-installer' \
  --set 'global.actions.execute=install' >"${RENDERED}"
contains 'serviceAccountName: existing-installer' || fail "existing installer account is not assigned"
pass "Existing installer ServiceAccount assignment"

helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'serviceAccount.enabled=true' \
  --set 'serviceAccount.create=true' \
  --set 'global.actions.execute=deploy' \
  --set 'global.tier[0].custom.serviceAccountName=legacy-tier-account' >"${RENDERED}"
contains 'serviceAccountName: legacy-tier-account' || fail "legacy tier override lost"
pass "Legacy tier override precedence"

expect_failure "create without enabled" \
  --set 'serviceAccount.create=true'

expect_failure "enabled without name or create" \
  --set 'serviceAccount.enabled=true'

expect_failure "invalid managed ServiceAccount name" \
  --set 'serviceAccount.enabled=true' \
  --set 'serviceAccount.create=true' \
  --set 'serviceAccount.name=Invalid_Name'

printf 'All Pega ServiceAccount checks passed.\n'
