#!/usr/bin/env bash
RELEASE_NAME=my-mysql

#helm install ${RELEASE_NAME} bitnami/mysql -f values.yaml
helm install ${RELEASE_NAME} \
--set auth.rootPassword=secretpassword,auth.database=app_database \
bitnami/mysql
