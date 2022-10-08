#!/usr/bin/env bash

curl -X POST -H "Content-Type: application/json" --data @message-mongo-sink.json http://localhost:8083/connectors -w "\n"
#curl -X POST -H "Content-Type: application/json" -d @message-mongo-source.json http://localhost:8083/connectors -w "\n"
# print all connectors added to kafka connect
curl -X GET http://localhost:8083/connectors

