#!/usr/bin/env bash

# Test structure:
# 1. Render the chart with ServiceAccount management disabled or enabled.
# 2. Assert generated accounts, workload references, token settings, and
#    precedence between global and tier-specific overrides.
# 3. Exercise validation failures for invalid ServiceAccount configuration.
#
# Test setup:
# - The chart is rendered with `helm template`; no Kubernetes resources are
#   created, updated, or deleted.
# - `global.provider=k8s` makes the chart pass provider validation.
# - Temporary rendered manifests and Helm errors are written under TMPDIR and
#   removed on exit.
# - Each case supplies its own values so cases do not depend on one another.

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
CHART="${SCRIPT_DIR}/../charts/pega"
RENDERED="${TMPDIR:-/tmp}/pega-serviceaccount-test.$$"
ERRORS="${RENDERED}.err"
trap 'rm -f "${RENDERED}" "${ERRORS}"' EXIT

pass() {
  printf 'PASS: %s\n' "$1"
}

case_header() {
  printf '\nCASE: %s\n' "$1"
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

# GIVEN ServiceAccount management is not configured
# WHEN the chart is rendered with default values
# THEN no managed ServiceAccount resource is emitted.
case_header "GIVEN disabled management WHEN rendering defaults THEN no ServiceAccount is created"
helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' >"${RENDERED}"
if contains '^kind: ServiceAccount$'; then
  fail "ServiceAccount rendered while management is disabled"
fi
pass "ServiceAccount management is disabled by default"

# GIVEN runtime and installer ServiceAccount management is enabled
# WHEN the chart is rendered for a deployment
# THEN both accounts are created and token automount remains disabled.
case_header "GIVEN managed accounts WHEN rendering a deployment THEN both accounts are created"
helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'global.serviceAccount.enabled=true' \
  --set 'global.serviceAccount.create=true' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=true' \
  --set 'global.actions.execute=deploy' >"${RENDERED}"
contains 'name: pega-serviceaccount' || fail "runtime ServiceAccount is missing"
contains 'name: pega-installer-serviceaccount' || fail "installer ServiceAccount is missing"
[[ "$(grep -Ec '^automountServiceAccountToken: false$' "${RENDERED}")" -eq 2 ]] \
  || fail "token automount is not disabled on both accounts"
pass "Managed runtime and installer ServiceAccounts"

# GIVEN externally managed runtime and installer accounts
# WHEN create is disabled and explicit names are supplied
# THEN no ServiceAccount is created and the names are assigned to workloads.
case_header "GIVEN pre-existing accounts WHEN create is disabled THEN accounts are referenced but not created"
helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'global.serviceAccount.enabled=true' \
  --set 'global.serviceAccount.create=false' \
  --set 'global.serviceAccount.name=existing-runtime' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=false' \
  --set 'installer.serviceAccount.name=existing-installer' \
  --set 'global.actions.execute=deploy' >"${RENDERED}"
if contains '^kind: ServiceAccount$'; then
  fail "external ServiceAccounts were unexpectedly created"
fi
contains 'serviceAccountName: existing-runtime' || fail "existing runtime account is not assigned"
pass "Pre-existing ServiceAccounts"

# GIVEN an externally managed installer account
# WHEN an install action is rendered
# THEN the installer Job references the supplied account.
case_header "GIVEN an external installer account WHEN rendering install THEN the installer uses it"
helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=false' \
  --set 'installer.serviceAccount.name=existing-installer' \
  --set 'global.actions.execute=install' >"${RENDERED}"
contains 'serviceAccountName: existing-installer' || fail "existing installer account is not assigned"
pass "Existing installer ServiceAccount assignment"

# GIVEN the shared global account is enabled
# WHEN all workloads are rendered with a custom global name
# THEN the runtime and supported subcharts use that account.
case_header "GIVEN a global account WHEN rendering workloads THEN the global account is propagated"
helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'global.serviceAccount.enabled=true' \
  --set 'global.serviceAccount.create=true' \
  --set 'global.serviceAccount.name=shared-runtime' \
  --set 'installer.serviceAccount.enabled=true' \
  --set 'installer.serviceAccount.create=true' \
  --set 'global.actions.execute=deploy' >"${RENDERED}"
contains 'name: shared-runtime' || fail "global runtime ServiceAccount is missing"
[[ "$(grep -Ec 'serviceAccountName: shared-runtime$' "${RENDERED}")" -ge 2 ]] \
  || fail "global runtime account was not propagated to workloads"
pass "Global ServiceAccount propagation"

# GIVEN a per-tier ServiceAccount override
# WHEN the shared account is enabled
# THEN the explicit tier override still has precedence.
case_header "GIVEN a tier override WHEN the global account is enabled THEN the tier account wins"
helm template pega "${CHART}" --namespace pega \
  --set 'global.provider=k8s' \
  --set 'global.serviceAccount.enabled=true' \
  --set 'global.serviceAccount.create=true' \
  --set 'global.actions.execute=deploy' \
  --set 'global.tier[0].custom.serviceAccountName=legacy-tier-account' >"${RENDERED}"
contains 'serviceAccountName: legacy-tier-account' || fail "legacy tier override lost"
pass "Tier override precedence"

# GIVEN create is enabled without enabling management
# WHEN the chart is rendered
# THEN Helm fails instead of silently ignoring the invalid configuration.
case_header "GIVEN create without enabled WHEN rendering THEN validation fails"
expect_failure "create without enabled" \
  --set 'global.serviceAccount.create=true'

# GIVEN management is enabled without a name or create=true
# WHEN the chart is rendered
# THEN Helm fails because the account cannot be resolved.
case_header "GIVEN enabled without name or create WHEN rendering THEN validation fails"
expect_failure "enabled without name or create" \
  --set 'global.serviceAccount.enabled=true'

# GIVEN an invalid Kubernetes ServiceAccount name
# WHEN the chart is rendered
# THEN Helm fails with a name validation error.
case_header "GIVEN an invalid account name WHEN rendering THEN validation fails"
expect_failure "invalid managed ServiceAccount name" \
  --set 'global.serviceAccount.enabled=true' \
  --set 'global.serviceAccount.create=true' \
  --set 'global.serviceAccount.name=Invalid_Name'

printf 'All Pega ServiceAccount checks passed.\n'
