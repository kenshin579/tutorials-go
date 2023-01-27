#!/usr/bin/env bash

cat json/ex2.json | jq '.[] | select(.id == "423be8de-9c04-4f0e-8ff0-545a8cb175b4") | {name, country}'
