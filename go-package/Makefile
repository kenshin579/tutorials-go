GOPATH     := $(shell go env GOPATH)
VERSION	   := $(shell git rev-parse --short=4 HEAD)
BUILDTIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%S%z")
APP        := go-package

LDFLAGS	   := "-X github.com/kenshin579/tutorials-go/go-package/config.Version=${VERSION} \
-X github.com/kenshin579/tutorials-go/go-package/config.BuildTime=${BUILDTIME}"

.PHONY: clean
clean:
	rm -f bin/$(APP)
	rmdir bin

.PHONY: build
build:
	go build -v -ldflags ${LDFLAGS} -o bin/$(APP) main.go

.PHONE: run
run: build
	bin/$(APP)
