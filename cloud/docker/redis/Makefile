REDIS_CONTAINER_NAME := go-redis

# 비밀번호 인증이 필요하면 아래 주석을 해제한다.
# 주석 상태에서도 `make redis-create REDIS_PASSWORD=xxx`로 일회성 지정이 가능하다.
#REDIS_PASSWORD ?= mypassword

.PHONY: redis-create
redis-create:
	@docker run --name $(REDIS_CONTAINER_NAME) -d -p 6379:6379 \
		redis redis-server $(if $(REDIS_PASSWORD),--requirepass $(REDIS_PASSWORD))

.PHONY: redis-delete
redis-delete:
	@docker stop $(REDIS_CONTAINER_NAME)
	@docker rm $(REDIS_CONTAINER_NAME)

redis-recreate: redis-delete redis-create

.PHONY: redis-cli
redis-cli:
	@docker exec -it $(REDIS_CONTAINER_NAME) redis-cli $(if $(REDIS_PASSWORD),-a $(REDIS_PASSWORD))
