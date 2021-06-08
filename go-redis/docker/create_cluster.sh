#!/usr/bin/env bash

echo "create cluster..."
docker exec -it redis7001 redis-cli -p 7001 -a password --cluster create 192.168.0.19:7001 192.168.0.19:7002 192.168.0.19:7003 192.168.0.19:7004 192.168.0.19:7005 192.168.0.19:7006 --cluster-replicas 1

