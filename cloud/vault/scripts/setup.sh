#!/usr/bin/env bash
# dev Vault에 Transit KEK, policy, policy별 token을 만든다. 여러 번 실행해도 결과가 같다.
set -euo pipefail

cd "$(dirname "$0")/.."

vault() {
  docker compose exec -T \
    -e VAULT_ADDR=http://127.0.0.1:8200 \
    -e VAULT_TOKEN=root \
    vault vault "$@"
}

if ! vault secrets list -format=json | grep -q '"transit/"'; then
  vault secrets enable transit
fi

vault write -f transit/keys/member-api-kek > /dev/null
echo "Transit KEK ready: member-api-kek"

vault policy write member-api - < policies/member-api.hcl
vault policy write member-api-migration - < policies/member-api-migration.hcl

mkdir -p .tokens
vault token create -policy=member-api -ttl=24h -field=token > .tokens/app.token
vault token create -policy=member-api-migration -ttl=24h -field=token > .tokens/migration.token

echo "setup done: .tokens/app.token, .tokens/migration.token"
