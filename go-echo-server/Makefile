REGISTRY 	:= kenshin579
APP    		:= advenoh
TAG         := go-echo-server
IMAGE       := $(REGISTRY)/$(APP):$(TAG)


.PHONY: docker-build
docker-build:
	@docker build -t $(IMAGE) -f Dockerfile .

.PHONY: docker-push
docker-push: docker-build
	@docker push $(IMAGE)
