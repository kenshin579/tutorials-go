#!/usr/bin/env bash


#http://www.compciv.org/recipes/cli/jq-for-parsing-json/
echo ""
cat sample1.json | jq '.name'

# Nested Objects
cat sample1.json | jq '.name.first'
cat sample1.json | jq '.bio .birthday'

# Accessing objects in an array
cat sample2.json | jq '.data.values'

#[]를 사용하면 배열의 값만 얻을 수 있음
cat sample2.json | jq '.data.values[]'
cat sample2.json | jq '.data.values[].id'

# select 함수 사용
# 배열 중에 id==id1인 데이터만 출력
cat sample2.json | jq '.data.values[]' | jq 'select(.id == "id1")'

# jq로 array에서 특정 필드로 매칭되는 필드 값을 추출하는 방법
cat json/ex2.json | jq '.[] | select(.id == "423be8de-9c04-4f0e-8ff0-545a8cb175b4") | {name, country}'

#JSON 문자열에서 배열의 수
echo '[{"username":"user1"},{"username":"user2"}]' | jq '. | length'
