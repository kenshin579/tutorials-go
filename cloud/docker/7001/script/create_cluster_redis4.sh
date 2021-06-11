#!/usr/bin/env bash

echo "starting..."

source /root/.bashrc

# todo : 실행이 잘 안됨
/data/redis-trib.rb -a create --replicas 1 192.168.0.131:7001 192.168.0.131:7002 192.168.0.131:7003 192.168.0.131:7004 192.168.0.131:7005 192.168.0.131:7006
