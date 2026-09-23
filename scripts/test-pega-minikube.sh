#!/usr/bin/env bash

# Test structure:
# 1. Validate the local Minikube context and installation inputs.
# 2. Render the chart and assert the expected ServiceAccount and
#    NetworkPolicy resources before touching the cluster.
# 3. Install Pega with `install-deploy` into an isolated namespace.
# 4. Verify Helm, workloads, ServiceAccounts, and NetworkPolicies.
# 5. Test allowed and denied PostgreSQL and external egress paths.
# 6. Optionally test a listener running on the local host.
# 7. Remove the test release, namespace, and diagnostic Pods unless --keep is used.
#
# Test setup:
# - The script uses `helm template` and `helm install`; it is an integration
#   test and therefore changes the selected cluster unless --render-only is used.
# - At least one --values file is required. The first file should contain the
#   normal Pega installation values; later files can contain policy overrides.
# - The default namespace is intentionally different from the normal `pega`
#   namespace to prevent accidental changes to an existing installation.
# - A PostgreSQL Pod selector is required for the database connectivity checks.
# - A netshoot image is pulled by the cluster for network connectivity checks.

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
CHART="${REPO_ROOT}/charts/pega"
CONTEXT="minikube"
RELEASE="pega-integration-test"
NAMESPACE="pega-integration-test"
POSTGRES_SELECTOR="cnpg.io/cluster=pega-postgres"
POSTGRES_PORT="5432"
EXTERNAL_HOST="1.1.1.1"
EXTERNAL_ALLOWED_PORT="443"
EXTERNAL_BLOCKED_PORT="80"
HOST_PORT=""
TIMEOUT="20m"
KEEP=false
RENDER_ONLY=false
VALUES_FILES=()
CREATED_NAMESPACE=false
INSTALLED=false
TEST_PODS=()

usage() {
  cat <<'USAGE'
Usage:
  test-pega-minikube.sh --values FILE [options]

Required:
  --values FILE                 Values file; repeat for base values and overlays.

Options:
  --context NAME                Kubernetes context (default: minikube)
  --release NAME                Helm release name
  --namespace NAME              Isolated test namespace
  --postgres-selector SELECTOR  PostgreSQL Pod selector
  --postgres-port PORT          PostgreSQL destination port (default: 5432)
  --host-port PORT              Also test host.minikube.internal:PORT
  --timeout DURATION            Helm install timeout (default: 20m)
  --keep                        Keep release, namespace, and diagnostic Pods
  --render-only                 Render and validate only; do not contact the cluster
  -h, --help                    Show this help
USAGE
}

pass() {
  printf 'PASS: %s\n' "$1"
}

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || fail "$1 is required"
}

case_header() {
  printf '\nCASE: %s\n' "$1"
}

cleanup() {
  local exit_code=$?
  if [[ "${KEEP}" == true || "${RENDER_ONLY}" == true ]]; then
    exit "${exit_code}"
  fi

  if ((${#TEST_PODS[@]} > 0)); then
    kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" delete pod \
      "${TEST_PODS[@]}" --ignore-not-found >/dev/null 2>&1 || true
  fi
  if [[ "${INSTALLED}" == true ]]; then
    helm uninstall "${RELEASE}" \
      --kube-context "${CONTEXT}" \
      --namespace "${NAMESPACE}" >/dev/null 2>&1 || true
  fi
  if [[ "${CREATED_NAMESPACE}" == true ]]; then
    kubectl --context "${CONTEXT}" delete namespace "${NAMESPACE}" \
      --ignore-not-found >/dev/null 2>&1 || true
  fi
  exit "${exit_code}"
}
trap cleanup EXIT

while (($# > 0)); do
  case "$1" in
    --values)
      (($# >= 2)) || fail "--values requires a file"
      VALUES_FILES+=("$2")
      shift 2
      ;;
    --context)
      (($# >= 2)) || fail "--context requires a value"
      CONTEXT="$2"
      shift 2
      ;;
    --release)
      (($# >= 2)) || fail "--release requires a value"
      RELEASE="$2"
      shift 2
      ;;
    --namespace)
      (($# >= 2)) || fail "--namespace requires a value"
      NAMESPACE="$2"
      shift 2
      ;;
    --postgres-selector)
      (($# >= 2)) || fail "--postgres-selector requires a value"
      POSTGRES_SELECTOR="$2"
      shift 2
      ;;
    --postgres-port)
      (($# >= 2)) || fail "--postgres-port requires a value"
      POSTGRES_PORT="$2"
      shift 2
      ;;
    --host-port)
      (($# >= 2)) || fail "--host-port requires a value"
      HOST_PORT="$2"
      shift 2
      ;;
    --timeout)
      (($# >= 2)) || fail "--timeout requires a value"
      TIMEOUT="$2"
      shift 2
      ;;
    --keep)
      KEEP=true
      shift
      ;;
    --render-only)
      RENDER_ONLY=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage >&2
      fail "unknown argument: $1"
      ;;
  esac
done

((${#VALUES_FILES[@]} > 0)) || {
  usage >&2
  fail "at least one --values file is required"
}

[[ "${NAMESPACE}" != "pega" ]] \
  || fail "refusing to use the existing pega namespace; choose an isolated namespace"

for values_file in "${VALUES_FILES[@]}"; do
  [[ -f "${values_file}" ]] || fail "values file does not exist: ${values_file}"
done
HELM_VALUES_ARGS=()
for values_file in "${VALUES_FILES[@]}"; do
  HELM_VALUES_ARGS+=(--values "${values_file}")
done

require_command helm
require_command kubectl

# GIVEN a controlled set of Pega values
# WHEN the chart is rendered before installation
# THEN the expected managed accounts and zero-trust policies are present.
case_header "GIVEN installation values WHEN rendering THEN security resources are present"
RENDERED="$(mktemp)"
helm template "${RELEASE}" "${CHART}" \
  --namespace "${NAMESPACE}" \
  "${HELM_VALUES_ARGS[@]}" \
  --set global.provider=k8s \
  --set-string global.actions.execute=install-deploy >"${RENDERED}"
grep -q '^kind: ServiceAccount$' "${RENDERED}" \
  || fail "no ServiceAccount was rendered"
grep -q 'name: pega-serviceaccount' "${RENDERED}" \
  || fail "runtime ServiceAccount was not rendered"
grep -q 'name: pega-installer-serviceaccount' "${RENDERED}" \
  || fail "installer ServiceAccount was not rendered"
grep -q 'name: .*networkpolicy-default-deny' "${RENDERED}" \
  || fail "default-deny NetworkPolicy was not rendered"
grep -q 'name: .*networkpolicy-dns' "${RENDERED}" \
  || fail "DNS NetworkPolicy was not rendered"
rm -f "${RENDERED}"
pass "Installation manifest contains ServiceAccounts and NetworkPolicies"

if [[ "${RENDER_ONLY}" == true ]]; then
  printf '\nRender-only mode complete; no cluster resources were changed.\n'
  exit 0
fi

# GIVEN an isolated namespace and a valid Pega installation configuration
# WHEN Helm runs install-deploy
# THEN the release is installed and the Pega workloads become ready.
case_header "GIVEN an isolated namespace WHEN installing install-deploy THEN Pega becomes ready"
if kubectl --context "${CONTEXT}" get namespace "${NAMESPACE}" >/dev/null 2>&1; then
  fail "namespace already exists: ${NAMESPACE}; choose another isolated namespace"
fi
kubectl --context "${CONTEXT}" create namespace "${NAMESPACE}" >/dev/null
CREATED_NAMESPACE=true

helm install "${RELEASE}" "${CHART}" \
  --kube-context "${CONTEXT}" \
  --namespace "${NAMESPACE}" \
  "${HELM_VALUES_ARGS[@]}" \
  --set global.provider=k8s \
  --set-string global.actions.execute=install-deploy \
  --wait \
  --wait-for-jobs \
  --timeout "${TIMEOUT}"
INSTALLED=true
kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" wait \
  --for=condition=Available deployment \
  --all \
  --timeout="${TIMEOUT}"
pass "Pega install-deploy completed"

# GIVEN a completed installation
# WHEN Helm and Kubernetes resources are inspected
# THEN the release, ServiceAccounts, and NetworkPolicies match the contract.
case_header "GIVEN an installed release WHEN inspecting resources THEN security settings are applied"
helm status "${RELEASE}" --kube-context "${CONTEXT}" --namespace "${NAMESPACE}" \
  | grep -q 'STATUS: deployed' \
  || fail "Helm release is not deployed"
kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
  get serviceaccount pega-serviceaccount >/dev/null \
  || fail "runtime ServiceAccount is missing"
kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
  get serviceaccount pega-installer-serviceaccount >/dev/null \
  || fail "installer ServiceAccount is missing"
for policy in default-deny dns database kafka tiers; do
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
    get networkpolicy -o name \
    | grep -q -- "-networkpolicy-${policy}$" \
    || fail "NetworkPolicy is missing: ${policy}"
done
[[ "$(kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
  get serviceaccount pega-serviceaccount \
  -o jsonpath='{.automountServiceAccountToken}')" == "false" ]] \
  || fail "runtime ServiceAccount token automount is not disabled"
pass "Release, ServiceAccounts, and built-in NetworkPolicies are present"

# GIVEN a ready Pega installation
# WHEN Pod identities and readiness are inspected
# THEN every Pega runtime Pod is ready and uses the managed runtime account.
case_header "GIVEN ready Pega workloads WHEN inspecting Pod identities THEN the runtime account is used"
mapfile -t PEGA_PODS < <(
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
    get pods -l component=Pega -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}'
)
((${#PEGA_PODS[@]} > 0)) || fail "no Pega runtime Pods were found"
for pod in "${PEGA_PODS[@]}"; do
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
    wait --for=condition=Ready "pod/${pod}" --timeout="${TIMEOUT}"
  [[ "$(kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
    get pod "${pod}" -o jsonpath='{.spec.serviceAccountName}')" == "pega-serviceaccount" ]] \
    || fail "Pod ${pod} does not use pega-serviceaccount"
done
pass "Pega runtime Pods are ready and use pega-serviceaccount"

create_test_pod() {
  local name="$1"
  local label="$2"
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" run "${name}" \
    --image=nicolaka/netshoot \
    --restart=Never \
    --labels="${label}" \
    --command -- sleep 3600 >/dev/null
  TEST_PODS+=("${name}")
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" wait \
    --for=condition=Ready "pod/${name}" --timeout=120s
}

# GIVEN a labeled diagnostic Pod and an unprivileged diagnostic Pod
# WHEN they connect to the configured PostgreSQL Pod
# THEN the labeled client is allowed on the database port and the other is denied.
case_header "GIVEN labeled and unlabeled clients WHEN testing PostgreSQL THEN policy selects only the labeled client"
create_test_pod "custom-policy-client" "network-test=custom-policy-client"
create_test_pod "unprivileged-client" "network-test=unprivileged-client"
POSTGRES_IP="$(
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" \
    get pod -l "${POSTGRES_SELECTOR}" -o jsonpath='{.items[0].status.podIP}'
)"
[[ -n "${POSTGRES_IP}" ]] || fail "no PostgreSQL Pod matched: ${POSTGRES_SELECTOR}"
kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec custom-policy-client -- \
  nc -zvw 5 "${POSTGRES_IP}" "${POSTGRES_PORT}" >/dev/null \
  || fail "labeled client could not reach PostgreSQL"
if kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec unprivileged-client -- \
    nc -zvw 5 "${POSTGRES_IP}" "${POSTGRES_PORT}" >/dev/null 2>&1; then
  fail "unprivileged client unexpectedly reached PostgreSQL"
fi
pass "PostgreSQL access is allowed only for the labeled client"

# GIVEN a labeled diagnostic Pod and a public endpoint
# WHEN the Pod connects to allowed and disallowed ports
# THEN only the configured external port is reachable.
case_header "GIVEN an external endpoint WHEN testing allowed and denied ports THEN only HTTPS is reachable"
kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec custom-policy-client -- \
  nc -zvw 5 "${EXTERNAL_HOST}" "${EXTERNAL_ALLOWED_PORT}" >/dev/null \
  || fail "labeled client could not reach ${EXTERNAL_HOST}:${EXTERNAL_ALLOWED_PORT}"
if kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec custom-policy-client -- \
    nc -zvw 5 "${EXTERNAL_HOST}" "${EXTERNAL_BLOCKED_PORT}" >/dev/null 2>&1; then
  fail "labeled client unexpectedly reached ${EXTERNAL_HOST}:${EXTERNAL_BLOCKED_PORT}"
fi
if kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec unprivileged-client -- \
    nc -zvw 5 "${EXTERNAL_HOST}" "${EXTERNAL_ALLOWED_PORT}" >/dev/null 2>&1; then
  fail "unprivileged client unexpectedly reached external HTTPS"
fi
pass "External egress is limited to the configured client and port"

if [[ -n "${HOST_PORT}" ]]; then
  # GIVEN a listener on the local host at host.minikube.internal:HOST_PORT
  # WHEN the labeled and unprivileged clients connect
  # THEN only the labeled client can reach the local terminal.
  case_header "GIVEN a local host listener WHEN testing host egress THEN only the labeled client connects"
  kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec custom-policy-client -- \
    nc -zvw 5 host.minikube.internal "${HOST_PORT}" >/dev/null \
    || fail "labeled client could not reach host.minikube.internal:${HOST_PORT}"
  if kubectl --context "${CONTEXT}" --namespace "${NAMESPACE}" exec unprivileged-client -- \
      nc -zvw 5 host.minikube.internal "${HOST_PORT}" >/dev/null 2>&1; then
    fail "unprivileged client unexpectedly reached the local host"
  fi
  pass "Local host egress is limited to the labeled client"
fi

printf '\nAll Pega Minikube integration checks passed.\n'
