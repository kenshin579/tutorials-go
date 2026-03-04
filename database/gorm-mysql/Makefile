.PHONY: up down reset test

up:
	docker compose up -d

down:
	docker compose down

reset:
	docker compose down -v && docker compose up -d

test:
	go test -v -count=1 -p 1 ./...
