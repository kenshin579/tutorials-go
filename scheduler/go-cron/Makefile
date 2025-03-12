REDIS_CONTAINER_NAME := go-redis

.PHONY: redis-create
redis-create:
	@docker run --name $(REDIS_CONTAINER_NAME) -d -p 6379:6379 redis

.PHONY: redis-delete
redis-delete:
	@docker stop $(REDIS_CONTAINER_NAME)
	@docker rm $(REDIS_CONTAINER_NAME)

redis-recreate: redis-delete redis-create
