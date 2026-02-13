.PHONY: run-livekit stop-livekit run-be run-fe install-fe

run-livekit:
	docker run -d --name livekit-server --rm \
		-p 7880:7880 \
		-p 7881:7881 \
		-p 7882:7882/udp \
		-e LIVEKIT_KEYS="devkey: secret" \
		livekit/livekit-server \
		--dev \
		--bind 0.0.0.0

stop-livekit:
	docker stop livekit-server

run-be:
	cd backend && go run main.go

run-fe:
	cd frontend && npm run dev

install-fe:
	cd frontend && npm install
