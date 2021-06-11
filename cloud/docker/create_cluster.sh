#!/usr/bin/env bash

echo "create cluster..."
# redis 6
#docker exec -it redis7001 redis-cli -p 7001 -a password --cluster create 192.168.0.19:7001 192.168.0.19:7002 192.168.0.19:7003 192.168.0.19:7004 192.168.0.19:7005 192.168.0.19:7006 --cluster-replicas 1

# redis 4
docker exec -it redis7001 apt-get update
docker exec -it redis7001 apt-get install -y wget
docker exec -it redis7001 wget https://raw.githubusercontent.com/redis/redis/4.0/src/redis-trib.rb
docker exec -it redis7001 chmod 0755 ./redis-trib.rb

#todo : ruby를 설치하는 코드가 필요함
docker exec -it redis7001 ./redis-trib.rb create --cluster-replicas 1 192.168.0.19:7001 192.168.0.19:7002 192.168.0.19:7003 192.168.0.19:7004 192.168.0.19:7005 192.168.0.19:7006
