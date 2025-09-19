#!/usr/bin/env bash

# Cluster Deployment Script
# Created: 2025-01-26
# Purpose: Centralized deployment management for Kubernetes cluster resources

# Strict error handling
set -euo pipefail

# Directories
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly NAMESPACE_DIR="${SCRIPT_DIR}/namespaces"
readonly HELM_CHARTS_DIR="${SCRIPT_DIR}/helm-charts"
readonly CUSTOM_VALUES_DIR="${HELM_CHARTS_DIR}/custom-values"

# Logging and output formatting
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_YELLOW='\033[1;33m'
readonly COLOR_RESET='\033[0m'

# Logging functions
log_info() {
  printf "${COLOR_GREEN}[INFO]${COLOR_RESET} %s\n" "$*" >&2
}

log_warning() {
  printf "${COLOR_YELLOW}[WARNING]${COLOR_RESET} %s\n" "$*" >&2
}

log_error() {
  printf "${COLOR_RED}[ERROR]${COLOR_RESET} %s\n" "$*" >&2
}

# Dependency checking
ensure_command_exists() {
  local cmd="$1"
  if ! command -v "$cmd" &>/dev/null; then
    log_error "Required command '$cmd' not found. Please install it."
    exit 1
  fi
}

# Validate cluster connectivity
validate_cluster_connection() {
  if ! kubectl cluster-info &>/dev/null; then
    log_error "Unable to connect to Kubernetes cluster. Check your kubeconfig."
    exit 1
  fi
  log_info "Kubernetes cluster connection verified"
}

# Namespace deployment
deploy_namespaces() {
  log_info "Deploying Kubernetes namespaces..."

  if [[ ! -d "${NAMESPACE_DIR}" ]]; then
    log_warning "Namespaces directory not found: ${NAMESPACE_DIR}"
    return
  fi

  kubectl apply -f "${NAMESPACE_DIR}/"
  log_info "Namespaces deployed successfully"
}

# Helm chart deployment
deploy_helm_charts() {
  log_info "Deploying Helm charts..."

  # Prometheus Stack
  deploy_prometheus_stack

  # Add more helm chart deployments here in the future
}

# Deploy Prometheus Stack
deploy_prometheus_stack() {
  local chart_name="kube-prometheus-stack"
  local namespace="monitoring"
  local values_file="${CUSTOM_VALUES_DIR}/${chart_name}/values.yaml"

  log_info "Deploying ${chart_name}..."

  if [[ ! -f "${values_file}" ]]; then
    log_warning "Custom values file not found for ${chart_name}. Using default values."
    values_file=""
  fi

  helm upgrade --install prometheus prometheus-community/"${chart_name}" \
    --namespace "${namespace}" \
    --create-namespace \
    ${values_file:+--values "${values_file}"} \
    --wait

  log_info "${chart_name} deployed successfully"
}

# Main deployment function
main() {
  log_info "Starting cluster deployment..."

  # Validate prerequisites
  ensure_command_exists "kubectl"
  ensure_command_exists "helm"
  validate_cluster_connection

  # Deploy resources
  deploy_namespaces
  deploy_helm_charts

  log_info "Cluster deployment completed successfully 🚀"
}

# Execute main function if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
