REDIS_CONTAINER_NAME := go-redis
REDIS_IMAGE := redis:7
REDIS_PASSWORD ?= mypassword

.PHONY: redis-create
redis-create:
	@docker run --name $(REDIS_CONTAINER_NAME) -d -p 6379:6379 \
		$(REDIS_IMAGE) redis-server --requirepass $(REDIS_PASSWORD)

.PHONY: redis-delete
redis-delete:
	@docker stop $(REDIS_CONTAINER_NAME)
	@docker rm $(REDIS_CONTAINER_NAME)

redis-recreate: redis-delete redis-create

.PHONY: redis-cli
redis-cli:
	@docker exec -it $(REDIS_CONTAINER_NAME) redis-cli -a $(REDIS_PASSWORD)
