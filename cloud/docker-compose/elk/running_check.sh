#!/usr/bin/env bash

echo "docker-compose"
curl -X GET "localhost:9200/_cat/nodes?v=true&pretty"

