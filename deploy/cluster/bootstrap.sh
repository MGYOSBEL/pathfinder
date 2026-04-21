#!/usr/bin/env bash
# bootstrap.sh - Bootstrap a Pathfinder kind cluster with ArgoCD
#
# Usage:
#   ./deploy/cluster/bootstrap.sh <env>
#   ./deploy/cluster/bootstrap.sh local
#   ./deploy/cluster/bootstrap.sh prod
#
# Prerequisites per environment:
#   deploy/cluster/overlays/<env>/secrets.env  (copy from secrets.env.example)
#
# What this script does:
#   1. Creates a kind cluster for the given environment
#   2. Installs ArgoCD via Helm with env values + secret values (admin password)
#   3. Applies a single bootstrap Application that points to this repo
#   4. ArgoCD takes over and syncs everything else (projects, applicationsets, apps)
#   5. Starts cloud-provider-kind for LoadBalancer support (macOS/Kind)

set -euo pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly VALID_ENVS=("local" "prod")

readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_YELLOW='\033[1;33m'
readonly COLOR_RESET='\033[0m'

log_info()    { printf "${COLOR_GREEN}[INFO]${COLOR_RESET} %s\n" "$*" >&2; }
log_warning() { printf "${COLOR_YELLOW}[WARN]${COLOR_RESET} %s\n" "$*" >&2; }
log_error()   { printf "${COLOR_RED}[ERROR]${COLOR_RESET} %s\n" "$*" >&2; }

usage() {
  echo "Usage: $(basename "$0") <env>"
  echo "  env: one of ${VALID_ENVS[*]}"
  exit 1
}

validate_env() {
  local env="$1"
  for valid in "${VALID_ENVS[@]}"; do
    [[ "$env" == "$valid" ]] && return 0
  done
  log_error "Invalid environment '${env}'. Valid: ${VALID_ENVS[*]}"
  exit 1
}

ensure_command() {
  local cmd="$1"
  if ! command -v "$cmd" &>/dev/null; then
    log_error "Required command '${cmd}' not found. Please install it."
    exit 1
  fi
}

load_secrets() {
  local env="$1"
  local secrets_file="${SCRIPT_DIR}/overlays/${env}/secrets.env"

  if [[ ! -f "$secrets_file" ]]; then
    log_error "Secrets file not found: ${secrets_file}"
    log_error "Copy ${secrets_file}.example to ${secrets_file} and fill in your values."
    exit 1
  fi

  # shellcheck disable=SC1090
  source "$secrets_file"

  if [[ -z "${ARGOCD_ADMIN_PASSWORD:-}" ]]; then
    log_error "ARGOCD_ADMIN_PASSWORD is not set in ${secrets_file}"
    exit 1
  fi
}

bcrypt_hash() {
  docker run --rm httpd:alpine htpasswd -bnBC 10 "" "$1" | tr -d ':\n'
}

create_cluster() {
  local env="$1"
  local config="${SCRIPT_DIR}/overlays/${env}/kind-config.yaml"
  local cluster_name="pathfinder-${env}"

  if [[ ! -f "$config" ]]; then
    log_error "kind config not found: ${config}"
    exit 1
  fi

  if kind get clusters 2>/dev/null | grep -q "^${cluster_name}$"; then
    log_warning "Cluster '${cluster_name}' already exists, skipping creation"
  else
    log_info "Creating kind cluster '${cluster_name}'..."
    kind create cluster --config "$config"
    log_info "Cluster '${cluster_name}' created"
  fi
}

install_argocd() {
  local env="$1"
  local password_hash="$2"
  local env_values="${SCRIPT_DIR}/overlays/${env}/argocd-values.yaml"

  log_info "Installing ArgoCD..."
  helm repo add argo https://argoproj.github.io/argo-helm --force-update
  helm repo update argo

  helm upgrade --install argocd argo/argo-cd \
    --namespace argocd \
    --create-namespace \
    --values "${env_values}" \
    --set-string "configs.secret.argocdServerAdminPassword=${password_hash}" \
    --wait

  log_info "ArgoCD installed"
}

apply_bootstrap_app() {
  local manifest="${SCRIPT_DIR}/argocd/bootstrap.yaml"

  if [[ ! -f "$manifest" ]]; then
    log_error "Bootstrap manifest not found: ${manifest}"
    exit 1
  fi

  log_info "Applying bootstrap Application..."
  kubectl apply -f "$manifest"
  log_info "Bootstrap Application applied — ArgoCD will sync the rest"
}

start_cloud_provider() {
  local env="$1"
  if [[ "$env" != "local" ]]; then
    return
  fi

  if ! command -v cloud-provider-kind &>/dev/null; then
    log_warning "Skipping cloud-provider-kind (not installed)"
    return
  fi

  # Kill any existing instance for this cluster
  sudo pkill -f "cloud-provider-kind" 2>/dev/null || true
  sleep 1

  log_info "Starting cloud-provider-kind in the background (requires sudo)..."
  sudo cloud-provider-kind > /tmp/cloud-provider-kind.log 2>&1 &
  log_info "cloud-provider-kind started (PID $!, log: /tmp/cloud-provider-kind.log)"
}

print_access_info() {
  echo ""
  log_info "Bootstrap complete 🚀"
  echo ""
  echo "  ArgoCD is now syncing deploy/cluster/argocd/ from the repo."
  echo ""
  echo "  Once ArgoCD finishes syncing all applications (cert-manager, envoy-gateway,"
  echo "  gateway-config), services will be available via the Gateway."
  echo ""
  echo "  Add to /etc/hosts (use the Gateway LoadBalancer IP once assigned):"
  echo "    127.0.0.1 grafana.pathfinder.local prometheus.pathfinder.local argocd.pathfinder.local jaeger.pathfinder.local alertmanager.pathfinder.local"
  echo ""
  echo "  Services (HTTPS — accept self-signed cert in browser):"
  echo "    ArgoCD:       https://argocd.pathfinder.local"
  echo "    Grafana:      https://grafana.pathfinder.local"
  echo "    Prometheus:    https://prometheus.pathfinder.local"
  echo "    Jaeger:        https://jaeger.pathfinder.local"
  echo "    Alertmanager:  https://alertmanager.pathfinder.local"
  echo ""
  echo "  Credentials:"
  echo "    ArgoCD — admin / (as set in overlays/<env>/secrets.env)"
  echo "    Grafana — admin / prom-operator (default from kube-prometheus-stack)"
  echo ""
  echo "  Useful commands:"
  echo "    kubectl get gateway -n networking          # check Gateway status"
  echo "    kubectl get httproute -n networking         # check route status"
  echo "    kubectl get svc -n envoy-gateway-system     # check LoadBalancer IP"
  echo ""
  echo "  Fallback (if Gateway is not yet ready):"
  echo "    kubectl port-forward svc/argocd-server -n argocd 8080:80"
  echo ""
}

main() {
  [[ $# -lt 1 ]] && usage

  local env="$1"

  validate_env "$env"
  ensure_command "kind"
  ensure_command "kubectl"
  ensure_command "helm"
  ensure_command "docker"

  if ! command -v cloud-provider-kind &>/dev/null; then
    log_warning "cloud-provider-kind not found. LoadBalancer services won't get external IPs."
    log_warning "Install it: go install sigs.k8s.io/cloud-provider-kind@latest"
  fi

  log_info "Bootstrapping Pathfinder cluster for environment: ${env}"

  load_secrets "$env"

  local password_hash
  password_hash=$(bcrypt_hash "$ARGOCD_ADMIN_PASSWORD")

  create_cluster "$env"
  install_argocd "$env" "$password_hash"
  apply_bootstrap_app
  start_cloud_provider "$env"
  print_access_info
}

main "$@"