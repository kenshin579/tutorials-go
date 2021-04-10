
.PHONY: clean
clean:
	-rm go.sum

.PHONY: package
package: clean
	go mod download
