#!/usr/bin/env bash
# bootstrap.sh - Bootstrap a Pathfinder kind cluster with ArgoCD
#
# Usage:
#   ./deploy/cluster/bootstrap.sh <env>
#   ./deploy/cluster/bootstrap.sh dev
#   ./deploy/cluster/bootstrap.sh staging
#   ./deploy/cluster/bootstrap.sh prod
#
# Prerequisites per environment:
#   deploy/cluster/environments/<env>/secrets.env  (copy from secrets.env.example)
#
# What this script does:
#   1. Creates a kind cluster for the given environment
#   2. Installs ArgoCD via Helm with env values + secret values (admin password)
#   3. Applies a single bootstrap Application that points to this repo
#   4. ArgoCD takes over and syncs everything else (projects, applicationsets, apps)

set -euo pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly VALID_ENVS=("dev" "staging" "prod")

# Temp file for secret values — global so the EXIT trap can clean it up
_SECRET_VALUES_TMP=""

cleanup() {
  [[ -n "${_SECRET_VALUES_TMP:-}" ]] && rm -f "$_SECRET_VALUES_TMP"
}
trap cleanup EXIT

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
  local secrets_file="${SCRIPT_DIR}/environments/${env}/secrets.env"

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
  # htpasswd is part of apache2-utils / httpd-tools
  htpasswd -bnBC 10 "" "$1" | tr -d ':\n'
}

create_cluster() {
  local env="$1"
  local config="${SCRIPT_DIR}/environments/${env}/kind-config.yaml"
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
  local env_values="${SCRIPT_DIR}/environments/${env}/values/argocd.yaml"

  # Write secret values to a temp file — cleaned up on EXIT, never persisted
  _SECRET_VALUES_TMP=$(mktemp)

  cat > "$_SECRET_VALUES_TMP" <<EOF
configs:
  secret:
    argocdServerAdminPassword: "${password_hash}"
EOF

  log_info "Installing ArgoCD..."
  helm repo add argo https://argoproj.github.io/argo-helm --force-update
  helm repo update argo

  helm upgrade --install argocd argo/argo-cd \
    --namespace argocd \
    --create-namespace \
    --values "${env_values}" \
    --values "${_SECRET_VALUES_TMP}" \
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

print_access_info() {
  echo ""
  log_info "Bootstrap complete 🚀"
  echo ""
  echo "  ArgoCD is now syncing deploy/cluster/argocd/ from the repo."
  echo ""
  echo "  ArgoCD UI:  https://localhost:8080"
  echo "  Username:   admin"
  echo "  Password:   (as set in environments/<env>/secrets.env)"
  echo ""
  echo "  To open the UI:"
  echo "    kubectl port-forward svc/argocd-server -n argocd 8080:80   # dev (insecure)"
  echo "    kubectl port-forward svc/argocd-server -n argocd 8080:443  # staging/prod (TLS)"
  echo ""
}

main() {
  [[ $# -lt 1 ]] && usage

  local env="$1"

  validate_env "$env"
  ensure_command "kind"
  ensure_command "kubectl"
  ensure_command "helm"
  ensure_command "htpasswd"

  log_info "Bootstrapping Pathfinder cluster for environment: ${env}"

  load_secrets "$env"

  local password_hash
  password_hash=$(bcrypt_hash "$ARGOCD_ADMIN_PASSWORD")

  create_cluster "$env"
  install_argocd "$env" "$password_hash"
  apply_bootstrap_app
  print_access_info
}

main "$@"