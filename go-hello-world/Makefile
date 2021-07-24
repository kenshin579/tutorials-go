DOCKERID := kenshin579
REGISTRY    := advenoh
TAG         := multi-arch-hello
IMAGE       := $(DOCKERID)/$(REGISTRY):$(TAG)

.PHONY: docker-push
docker-push:
	@docker buildx build \
	--platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 \
	-t $(IMAGE) --push -f Dockerfile .
