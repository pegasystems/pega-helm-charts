#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
CHART="${REPO_ROOT}/charts/pega"
NAMESPACE="networkpolicy-test"
RENDERED="${TMPDIR:-/tmp}/pega-networkpolicy-test.$$"

cleanup() {
  rm -f "${RENDERED}" "${RENDERED}.err"
}
trap cleanup EXIT

pass() {
  printf 'PASS: %s\n' "$1"
}

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

require_contains() {
  local file="$1"
  local pattern="$2"
  local description="$3"
  grep -Eq -- "${pattern}" "${file}" || fail "${description}"
}

require_not_contains() {
  local file="$1"
  local pattern="$2"
  local description="$3"
  ! grep -Eq -- "${pattern}" "${file}" || fail "${description}"
}

render_policy_templates() {
  local output="$1"
  shift
  helm template pega "${CHART}" \
    --namespace "${NAMESPACE}" \
    --set 'global.provider=k8s' \
    --set-json 'networkPolicy.customPolicies=[]' \
    "$@" > "${output}"
}

expect_render_failure() {
  local description="$1"
  shift
  if helm template pega "${CHART}" \
    --namespace "${NAMESPACE}" \
    --set 'global.provider=k8s' \
    "$@" > "${RENDERED}" 2> "${RENDERED}.err"; then
    cat "${RENDERED}.err" >&2
    fail "${description} unexpectedly rendered successfully"
  fi
  pass "${description}"
}

command -v helm >/dev/null 2>&1 || fail "helm is required"
helm version --short >/dev/null

helm lint "${CHART}" >/dev/null
pass "Helm lint"

if render_policy_templates "${RENDERED}" --set 'networkPolicy.enabled=false'; then
  require_not_contains "${RENDERED}" 'kind: NetworkPolicy' \
    "disabled NetworkPolicy feature emitted a policy"
  pass "NetworkPolicy generation is disabled by default"
fi

render_policy_templates "${RENDERED}" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.defaultDeny=true' \
  --set-string 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16'
require_contains "${RENDERED}" 'name: pega-networkpolicy-default-deny' \
  "default-deny policy was not rendered"
require_contains "${RENDERED}" 'name: pega-networkpolicy-dns' \
  "DNS policy was not rendered"
require_contains "${RENDERED}" 'kubernetes.io/metadata.name: kube-system' \
  "DNS namespace selector is missing"
require_contains "${RENDERED}" 'k8s-app: kube-dns' \
  "DNS pod selector is missing"
require_contains "${RENDERED}" 'cidr: "10.10.0.0/16"' \
  "configured ingress CIDR is missing"
require_contains "${RENDERED}" 'component: Pega' \
  "same-namespace Pega ingress is missing"
pass "Default-deny, DNS, and ingress policies"

render_policy_templates "${RENDERED}" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.defaultDeny=false' \
  --set-string 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16'
require_not_contains "${RENDERED}" 'name: pega-networkpolicy-default-deny' \
  "default-deny policy rendered when disabled"
require_contains "${RENDERED}" 'name: pega-networkpolicy-tiers' \
  "tier policy was not rendered with default-deny disabled"
pass "Default-deny can be disabled without disabling allow policies"

render_policy_templates "${RENDERED}" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.defaultDeny=true' \
  --set-string 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.database.enabled=true' \
  --set-string 'networkPolicy.database.cidrs[0]=10.20.30.40/32' \
  --set 'networkPolicy.database.ports[0]=5432' \
  --set 'networkPolicy.kafka.enabled=true' \
  --set-string 'networkPolicy.kafka.cidrs[0]=10.40.0.10/32' \
  --set 'networkPolicy.kafka.ports[0]=9092' \
  --set 'networkPolicy.externalSearch.enabled=true' \
  --set-string 'networkPolicy.externalSearch.cidrs[0]=10.50.0.0/24' \
  --set 'networkPolicy.externalSearch.ports[0]=443' \
  --set-string 'networkPolicy.installer.cidrs[0]=10.60.0.0/24' \
  --set 'networkPolicy.installer.ports[0]=443' \
  --set 'networkPolicy.customPolicies[0].name=allow-monitoring' \
  --set 'networkPolicy.customPolicies[0].podSelector.matchLabels.component=Pega' \
  --set 'networkPolicy.customPolicies[0].policyTypes[0]=Ingress' \
  --set-string 'networkPolicy.customPolicies[0].ingress[0].from[0].ipBlock.cidr=10.70.0.0/24' \
  --set 'networkPolicy.customPolicies[0].ingress[0].ports[0].protocol=TCP' \
  --set 'networkPolicy.customPolicies[0].ingress[0].ports[0].port=8080'

for policy in \
  default-deny dns tiers hazelcast cassandra installer-cassandra \
  database kafka installer-database external-search installer-additional \
  allow-monitoring; do
  require_contains "${RENDERED}" "name: pega-networkpolicy-${policy}" \
    "expected policy ${policy} was not rendered"
done

if [[ "$(grep -Ec '^  name: .*networkpolicy' "${RENDERED}")" -ne \
      "$(grep -E '^  name: .*networkpolicy' "${RENDERED}" | sort -u | wc -l | tr -d ' ')" ]]; then
  fail "rendered NetworkPolicy names are not unique"
fi
if grep -E '^  namespace:' "${RENDERED}" | grep -v "namespace: ${NAMESPACE}" >/dev/null; then
  fail "a rendered NetworkPolicy has an unexpected namespace"
fi
pass "Complete configuration, policy separation, unique names, and namespace scoping"

render_policy_templates "${RENDERED}" \
  --set 'networkPolicy.enabled=true' \
  --set-string 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.database.enabled=true' \
  --set 'networkPolicy.database.namespaceSelector.matchLabels.kubernetes\.io/metadata\.name=database' \
  --set 'networkPolicy.database.podSelector.matchLabels.role=primary' \
  --set 'networkPolicy.database.ports[0]=5432'
require_contains "${RENDERED}" 'kubernetes.io/metadata.name: database' \
  "database namespace selector is missing"
require_contains "${RENDERED}" 'role: primary' \
  "database pod selector is missing"
pass "Combined namespaceSelector and podSelector"

expect_render_failure "invalid database port" \
  --set 'networkPolicy.enabled=true' \
  --set-string 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.database.enabled=true' \
  --set-string 'networkPolicy.database.cidrs[0]=10.20.30.40/32' \
  --set 'networkPolicy.database.ports[0]=0'

LONG_NAME="$(printf 'a%.0s' {1..80})"
render_policy_templates "${RENDERED}" \
  --set 'networkPolicy.enabled=true' \
  --set-string 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set "global.deployment.name=${LONG_NAME}" \
  --set 'networkPolicy.customPolicies[0].name=custom-policy'
while read -r name; do
  [[ "${#name}" -le 63 ]] || fail "policy name exceeds 63 characters: ${name}"
done < <(grep -E '^  name: .*networkpolicy' "${RENDERED}" | awk '{print $2}')
pass "Long deployment names produce valid-length policy names"

expect_render_failure "missing ingress peer" \
  --set 'networkPolicy.enabled=true'
expect_render_failure "database peer is required" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.database.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.database.ports[0]=5432'
expect_render_failure "database port is required" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.database.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set-string 'networkPolicy.database.cidrs[0]=10.20.30.40/32'
expect_render_failure "Kafka peer is required" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.kafka.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.kafka.ports[0]=9092'
expect_render_failure "external Search port is required" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.externalSearch.enabled=true' \
  --set-string 'networkPolicy.externalSearch.cidrs[0]=10.50.0.0/24'
expect_render_failure "installer peer is required" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.installer.ports[0]=443'
expect_render_failure "invalid custom policy name" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.customPolicies[0].name=invalid_name'
expect_render_failure "duplicate custom policy name" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.customPolicies[0].name=duplicate' \
  --set 'networkPolicy.customPolicies[1].name=duplicate'
expect_render_failure "custom policy name colliding with built-in policy" \
  --set 'networkPolicy.enabled=true' \
  --set 'networkPolicy.ingress.cidrs[0]=10.10.0.0/16' \
  --set 'networkPolicy.customPolicies[0].name=tiers'

printf 'All Pega NetworkPolicy checks passed.\n'
