# 앱(member-api) 권한: wrap된 DEK를 읽고 KEK로 unwrap만 한다
path "secret/data/member-api/dek" {
  capabilities = ["read"]
}

path "transit/decrypt/member-api-kek" {
  capabilities = ["update"]
}
