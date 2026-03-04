.PHONY: generate server client clean test bench lint

generate:
	buf generate

lint:
	buf lint

server:
	go run ./cmd/server

client:
	go run ./cmd/client

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem -count=3 ./...

clean:
	rm -rf gen/
