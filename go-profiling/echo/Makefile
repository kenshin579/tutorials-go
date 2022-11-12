REGISTRY 	:= kenshin579
APP    		:= advenoh
TAG         := go-profile-server
IMAGE       := $(REGISTRY)/$(APP):$(TAG)
PROFILE_CONTAINER_NAME := go-profile-server

.PHONY: package
package:
	go mod tidy

.PHONY: docker-build
docker-build:
	@docker build -t $(IMAGE) -f Dockerfile .

.PHONY: docker-push
docker-push: docker-build
	@docker push $(IMAGE)

.PHONY: profile-create
profile-create:
	@docker run --name $(PROFILE_CONTAINER_NAME) -d $(IMAGE)

.PHONY: profile-delete
profile-delete:
	@docker stop $(PROFILE_CONTAINER_NAME)
	@docker rm $(PROFILE_CONTAINER_NAME)
