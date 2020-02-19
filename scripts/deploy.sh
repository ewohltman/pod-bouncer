#!/usr/bin/env bash

set -e
set -o pipefail # Only exit with zero if all commands of the pipeline exit successfully

SCRIPT_PATH=$(readlink -f "${0}")
SCRIPT_DIR=$(dirname "${SCRIPT_PATH}")

COMMIT=$(git rev-parse --short HEAD)

REPO_YMLS="${SCRIPT_DIR}/../deployments/kubernetes"

NAMESPACE_YML="${REPO_YMLS}/namespace.yml"
SERVICE_YML="${REPO_YMLS}/service.yml"

DEPLOYMENT_YML="${REPO_YMLS}/deployment.yml"
VARIABLIZED_DEPLOYMENT_YML="/tmp/deployment.yml"

setup() {
  cp "${DEPLOYMENT_YML}" "${VARIABLIZED_DEPLOYMENT_YML}"
}

applyValues() {
  sed -i "s|{COMMIT}|${COMMIT}|g" "${VARIABLIZED_DEPLOYMENT_YML}"
}

deploy() {
  kubectl apply -f "${NAMESPACE_YML}"
  kubectl apply -f "${SERVICE_YML}"
  kubectl apply -f "${VARIABLIZED_DEPLOYMENT_YML}"
  kubectl -n ephemeral-roles rollout status --timeout 120s deployment/pod-bouncer
}

cleanup() {
  rm -f "${VARIABLIZED_DEPLOYMENT_YML}"
}

trap cleanup EXIT

setup
applyValues
deploy
