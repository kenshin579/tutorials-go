REGISTRY 	:= kenshin579
APP         := dockertest
TAG         := latest
IMAGE       := $(REGISTRY)/$(APP):$(TAG)


.PHONY: docker-build
docker-build:
	@docker build -t $(IMAGE) -f Dockerfile .

.PHONY: docker-run
docker-run: docker-build
	@docker run $(IMAGE)

.PHONY: zip
zip:
	@zip eb.zip -r *

all:
	echo package, build, test

.PHONY: package
package:
	go mod tidy

.PHONY: build
build: package test
	go build -v -o $(APP) main.go

.PHONY: test
test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: clean
	go clean
	rm -rf bin
