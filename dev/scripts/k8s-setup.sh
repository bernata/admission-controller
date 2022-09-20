#!/bin/sh

YAML_DIR="$(dirname "$0")"/..
CERTS_DIR="$(dirname "$0")"/../certs

echo "Creating namespace ..."
kubectl --context kind-demo apply -f "${YAML_DIR}/namespace.yaml"

echo "Registering validating webhook"
# Replace {CA_PEM_BASE64} in yaml file with base64 encoded blob
CA_PEM_BASE64="$(base64 "${CERTS_DIR}/ca.pem")"
sed -e 's@${CA_PEM_BASE64}@'"$CA_PEM_BASE64"'@g' <"${YAML_DIR}/validating-webhook.yaml" \
    | kubectl --context kind-demo apply -f -
