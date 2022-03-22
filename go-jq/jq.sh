#!/usr/bin/env bash


#http://www.compciv.org/recipes/cli/jq-for-parsing-json/
echo ""
cat sample.json | jq '.name'

# Nested Objects
cat sample.json | jq '.name.first'
cat sample.json | jq '.bio .birthday'

# Accessing objects in an array
cat sample1.json | jq '.data.values'

#[]를 사용하면 배열의 값만 얻을 수 있음
cat sample1.json | jq '.data.values[]'
cat sample1.json | jq '.data.values[].id'

# select 함수 사용
# 배열 중에 id==id1인 데이터만 출력
cat sample1.json | jq '.data.values[]' | jq 'select(.id == "id1")'
