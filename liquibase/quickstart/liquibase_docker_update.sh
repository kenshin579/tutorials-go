#!/usr/bin/env bash


# docker network create mysql-network
# docker network connect mysql-network go-mysql

docker run --rm --network mysql-network \
  -v $(pwd):/liquibase/changelog \
  -e INSTALL_MYSQL=true \
  liquibase/liquibase \
  --log-level=info \
  --defaultsFile=/liquibase/changelog/liquibase.docker.properties update
