#!/usr/bin/env bash

echo "run docker private-registry"

# mkdir -p /Users/user/data/docker/auth
# 사용자 인증 파일 생성 (username: admin, password: password)
# htpasswd -Bbn admin password > /Users/user/data/docker/auth/htpasswd

docker run -d --name private-registry -p 7001:5000 \
  --restart=always \
  -v /Users/user/data/docker/private-registry:/var/lib/registry \
  -v /Users/user/data/docker/auth:/auth \
  -e "REGISTRY_AUTH=htpasswd" \
  -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
  -e "REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd" \
  registry:2

#--mount type=volume,source=registry-volume,target=/var/lib/registry \
#-v /home/evan/cert:/certs \
#-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/private_registry.crt \
#-e REGISTRY_HTTP_TLS_KEY=/certs/private_registry.key \
#-v $(pwd)/auth:/auth \

