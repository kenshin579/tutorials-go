
.PHONY: clean
	go clean
	rm -rf bin


.PHONY: package
package:
	go mod tidy

