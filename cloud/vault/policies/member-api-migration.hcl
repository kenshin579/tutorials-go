# migration 권한: DEK 발급·적재·교체와 재암호화에 필요한 만큼만 준다
# KEK 자체(transit/keys/*)는 다루지 않는다
path "transit/datakey/wrapped/member-api-kek" {
  capabilities = ["update"]
}

path "transit/decrypt/member-api-kek" {
  capabilities = ["update"]
}

path "secret/data/member-api/dek" {
  capabilities = ["create", "update", "read"]
}
