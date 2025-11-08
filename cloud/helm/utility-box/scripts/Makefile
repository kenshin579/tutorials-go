REGISTRY    := localhost:7001
APP         := utility-box
USERNAME    := username
PASSWORD    := password
PLATFORM	:= linux/arm64
IMAGE       := $(REGISTRY)/$(APP):v0.1

.PHONY: docker-build
docker-build:
	@echo $(IMAGE)
	@docker build --platform $(PLATFORM) -t $(IMAGE) -f Dockerfile .

.PHONY: docker-push
docker-push: docker-build
	@docker login --username $(USERNAME) --password $(PASSWORD) $(REGISTRY)
	@docker push $(IMAGE)
