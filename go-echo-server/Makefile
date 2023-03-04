REGISTRY 	:= kenshin579
APP    		:= go-echo-server
TAG         := v0.1
IMAGE       := $(REGISTRY)/$(APP):$(TAG)


.PHONY: docker-build
docker-build:
	@docker build -t $(IMAGE) -f Dockerfile .

.PHONY: docker-push
docker-push: docker-build
	@docker push $(IMAGE)

.PHONY: clean
	go clean
	rm -rf bin

.PHONY: package
package:
	go mod tidy

.PHONY: swagger
swagger:
	@go get -d github.com/swaggo/swag/cmd/swag@v1.8.7
	@go install github.com/swaggo/swag/cmd/swag@v1.8.7
	@swag i --parseDepth=3 --parseDependency -g cmd/server/main.go
	@go mod tidy
