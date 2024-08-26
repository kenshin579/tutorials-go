#!/usr/bin/env bash

#############################################################################################
# Initialize
#############################################################################################

set -e

#############################################################################################
# Constants
#############################################################################################

declare -A docker_image_map

# harbor에서 frank-test 프로젝트를 생성한다
docker_image_map['docker://docker.io/library/mariadb:latest']='docker://demo.goharbor.io/frank-test/mariadb:latest'

#############################################################################################
# Main
#############################################################################################

# login을 미리 해둔다
# skopeo login demo.goharbor.io

for src in "${!docker_image_map[@]}"; do
  echo "copying \"$src\" -> \"${docker_image_map[$src]}\""
  skopeo copy $src ${docker_image_map[$src]} --override-os linux --override-arch amd64
done

