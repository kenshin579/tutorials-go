#!/usr/bin/env bash

echo "installing sealed-secrets"

# Add the sealed-secrets helm repo
helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets

# Update the helm repositories
helm repo update

# Install or upgrade the sealed-secrets release
helm upgrade --install sealed-secrets sealed-secrets/sealed-secrets --namespace kube-system
