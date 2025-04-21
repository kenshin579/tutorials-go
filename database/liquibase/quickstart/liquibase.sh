#!/usr/bin/env bash

if [ $# -lt 1 ]; then
  echo "Usage: $0 <command>"
  exit 1
fi

COMMAND=$1
LIQUIBASE_PROPERTY=liquibase.properties

shift 1
EXTRA_ARGS="$@"
echo "COMMAND: $COMMAND | LIQUIBASE_PROPERTY: $LIQUIBASE_PROPERTY | EXTRA_ARGS: $EXTRA_ARGS"

# Liquibase 실행
case $COMMAND in
  update-one)
    liquibase --defaults-file=$PROPERTY_FILE update-count --count 1 $EXTRA_ARGS
    ;;
  update-all)
    liquibase --defaults-file=$PROPERTY_FILE update $EXTRA_ARGS
    ;;
  rollback-one)
    liquibase --defaults-file=$PROPERTY_FILE rollback-count --count 1 $EXTRA_ARGS
    ;;
  status)
    liquibase --defaults-file=$PROPERTY_FILE status $EXTRA_ARGS
    ;;
  history)
    liquibase --defaults-file=$PROPERTY_FILE history $EXTRA_ARGS
    ;;
  *)
    echo "unknown sub command"
    exit 1
    ;;
esac

