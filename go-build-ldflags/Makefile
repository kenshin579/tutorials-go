all:
	$(eval GIT_INFO=$(shell git show --pretty=format:%cs --no-patch)-$(shell git show --pretty=format:%h --no-patch))
	$(eval DATE=$(shell date))

	go build -ldflags="-X 'main.GitInfo=$(GIT_INFO)' -X 'main.BuildTime=$(DATE)'"

clean:
	go clean

test:
	go test
