#!/usr/bin/env bash

docker run -d \
  --name local-postgres \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=sampledb \
  -p 5432:5432 \
  postgres:16

