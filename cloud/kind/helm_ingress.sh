#!/usr/bin/env bash

if [ $# -lt 1 ]; then
  echo "Usage: $0 <command>"
  exit 1
fi

COMMAND=$1

echo "installing ingress-nginx"

case $COMMAND in
 install)
    helm install my-ingress-nginx ingress-nginx/ingress-nginx -f ingress-nginx-values.yaml -n ingress-nginx --create-namespace
    ;;
  uninstall)
    helm -n ingress-nginx  uninstall my-ingress-nginx
    ;;
  *)
    echo "unknown
    command"
    exit 1
    ;;
esac
