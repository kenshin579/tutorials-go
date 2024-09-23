#!/usr/bin/env bash


liquibase --defaultsFile=liquibase.properties \
  update
#  --log-level=debug \
#  --searchPath=db/changelog \

