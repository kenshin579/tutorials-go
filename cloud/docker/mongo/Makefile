MONGO_CONTAINER_NAME := go-mongo

.PHONY: mongo-create
mongo-create:
	@docker run --name $(MONGO_CONTAINER_NAME) -d \
	-e MONGO_INITDB_ROOT_USERNAME=mongo \
	-e MONGO_INITDB_ROOT_PASSWORD=mongopassword \
	-p 27017:27017 mongo

.PHONY: mongo-delete
mongo-delete:
	@docker stop $(MONGO_CONTAINER_NAME)
	@docker rm $(MONGO_CONTAINER_NAME)

.PHONY: mongo-recreate
mongo-recreate: mongo-delete mongo-create
