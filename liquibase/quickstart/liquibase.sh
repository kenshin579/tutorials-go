#!/usr/bin/env bash

if [ $# -lt 1 ]; then
  echo "Usage: $0 <command>"
  exit 1
fi

COMMAND=$1
LIQUIBASE_PROPERTY=liquibase.properties

#echo "COMMAND: $COMMAND | LIQUIBASE_PROPERTY: $LIQUIBASE_PROPERTY"

# Liquibase 실행
case $COMMAND in
  update-one)
    liquibase --defaults-file=$PROPERTY_FILE update-count --count 1
    ;;
  update-all)
    liquibase --defaults-file=$PROPERTY_FILE update
    ;;
  rollback-one)
    liquibase --defaults-file=$PROPERTY_FILE rollback-count --count 1
    ;;
  status)
    liquibase --defaults-file=$PROPERTY_FILE status
    ;;
  history)
    liquibase --defaults-file=$PROPERTY_FILE history
    ;;
  *)
    echo "unknown sub command"
    exit 1
    ;;
esac

