.PHONY: generate server client gateway clean test lint

generate:
	buf generate

lint:
	buf lint

server:
	go run ./cmd/server

client:
	go run ./cmd/client

gateway:
	go run ./cmd/gateway

test:
	go test -v ./...

clean:
	rm -rf gen/
