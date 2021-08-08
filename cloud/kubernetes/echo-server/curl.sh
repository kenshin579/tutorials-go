#!/usr/bin/env bash

SERVER=localhost:61535

# CUSTOM HTTP STATUS CODE
curl -I --header 'X-ECHO-CODE: 404' ${SERVER}
curl -I ${SERVER}/?echo_code=404

# Custom body
curl --header 'X-ECHO-BODY: amazing' localhost:8080
curl localhost:8080/?echo_body=amazing

# Custom response latency
curl --header 'X-ECHO-TIME: 5000' localhost:8080
curl "localhost:8080/?echo_time=5000"
