#!/usr/bin/env bash

echo "producer..."
kafka-console-producer.sh --bootstrap-server localhost:9092 --topic Tutorial2.pets < message.json
