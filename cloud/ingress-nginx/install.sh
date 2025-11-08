#!/usr/bin/env bash
kubectl get ns ingress-nginx &>/dev/null || kubectl create ns ingress-nginx

helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --set controller.replicaCount=2 \
  --set controller.resources.requests.cpu=100m \
  --set controller.resources.requests.memory=256Mi \
  --set controller.service.type=LoadBalancer
