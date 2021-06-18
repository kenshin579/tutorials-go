#!/usr/bin/env bash

kubectl delete -f volumes.yaml
#kubectl delete namespace redis-cluster-example
helm delete my-redis-cluster
