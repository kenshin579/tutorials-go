#!/usr/bin/env bash

helm delete my-jenkins
kubectl delete -f volumes.yaml

