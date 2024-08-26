#!/usr/bin/env bash

#############################################################################################
# Initialize
#############################################################################################
set -e

#############################################################################################
# Constants
#############################################################################################

DOCKER_REGISTRY_SRC_ADDR="docker.io"
DOCKER_REGISTRY_DST_ADDR="demo.goharbor.io"

# download dockker image
docker pull --platform=linux $DOCKER_REGISTRY_SRC_ADDR/library/mariadb:latest

# tag docker image
docker tag $DOCKER_REGISTRY_SRC_ADDR/library/mariadb:latest $DOCKER_REGISTRY_DST_ADDR/library/mariadb:latest

# login to docker registry
docker login $DOCKER_REGISTRY_DST_ADDR

# push docker image
docker push $DOCKER_REGISTRY_DST_ADDR/library/mariadb:latest
