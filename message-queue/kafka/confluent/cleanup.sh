#!/usr/bin/env bash

echo "cleanup..."
docker-compose stop
docker system prune -a --volumes --filter "label=io.confluent.docker"
