MONGO_CONTAINER_NAME := go-mongo
USERNAME=mongoadmin
PASSWORD=mongopassword

.PHONY: mongo-recreate
mongo-recreate: mongo-delete mongo-create

.PHONY: mongo-create
mongo-create:
	@docker run --name $(MONGO_CONTAINER_NAME) -d \
	-p 27017:27017 mongo:4.2

.PHONY: mongo-delete
mongo-delete:
	@docker stop $(MONGO_CONTAINER_NAME)
	@docker rm $(MONGO_CONTAINER_NAME)

