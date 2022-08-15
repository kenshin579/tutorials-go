#!/usr/bin/env bash

go run main.go -brokers localhost:29092 -group my_first_application -topics my_topic -verbose
