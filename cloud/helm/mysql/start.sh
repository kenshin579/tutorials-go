#!/usr/bin/env bash
RELEASE_NAME=my-mysql

helm install ${RELEASE_NAME} bitnami/mysql -f values.yaml
