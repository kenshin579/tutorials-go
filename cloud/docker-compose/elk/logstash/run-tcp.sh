#!/usr/bin/env bash

echo "running logstash"

# full path로 실행하지 않으면 conf 파일을 찾지 못하는 이슈가 있음
FULL_PATH=`readlink -f logstash-tcp.conf`

logstash -f $FULL_PATH

# run the following command from termainl
sleep 5
echo "hello message" | nc localhost 9600
