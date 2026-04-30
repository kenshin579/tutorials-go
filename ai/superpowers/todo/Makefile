.PHONY: dev-be dev-fe build preview-fe test test-be test-fe test-e2e install clean

BE_DIR := backend
FE_DIR := frontend

install:
	cd $(FE_DIR) && npm install

dev-be:
	cd $(BE_DIR) && go run .

dev-fe:
	cd $(FE_DIR) && npm run dev

build:
	cd $(FE_DIR) && npm run build

preview-fe: build
	cd $(FE_DIR) && npm run preview

test: test-be test-fe

test-be:
	cd $(BE_DIR) && go test -race ./... && go vet ./...

test-fe:
	cd $(FE_DIR) && npm test -- --run

test-e2e:
	cd $(FE_DIR) && npm run test:e2e

clean:
	rm -rf $(FE_DIR)/dist $(FE_DIR)/node_modules $(FE_DIR)/playwright-report $(FE_DIR)/test-results
