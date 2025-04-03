#!/usr/bin/env bash

export selector=productpage
export pod=`kubectl get pods -l app=${selector} -n default -o jsonpath='{.items[0].metadata.name}'`
export pod_ip=`kubectl get pod $pod -n default -o jsonpath='{.status.podIP}'`
curl -v -H "foo: bar" http://$pod_ip:9080/health
