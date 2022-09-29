#!/usr/bin/env bash

echo "consuming..."
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic messages.source
