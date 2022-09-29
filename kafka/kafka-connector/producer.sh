#!/usr/bin/env bash

echo "producer..."
kafka-console-producer.sh --bootstrap-server localhost:9092 --topic messages.source < message.json
