#!/usr/bin/env bash

echo "create cluster..."
# redis 6
docker exec -it redis7001 redis-cli -p 7001 -a password --cluster create 192.168.0.19:7001 192.168.0.19:7002 192.168.0.19:7003 192.168.0.19:7004 192.168.0.19:7005 192.168.0.19:7006 --cluster-replicas 1

# todo: redis 4 - redis-trib.rb로 cluster 생성하는 부분 잘 안됨
#docker exec -it redis7001 preinstall.sh
#docker exec -it redis7001 create_cluster_redis4.sh


