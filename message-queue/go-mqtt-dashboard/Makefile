.PHONY: help setup run-broker run-be run-fe stop-broker clean

# 기본 타겟
help:
	@echo "MQTT Dashboard - Available Commands"
	@echo "===================================="
	@echo "make setup        - 프로젝트 초기 설정 (의존성 설치)"
	@echo "make run-broker   - Mosquitto 브로커 실행"
	@echo "make run-be       - Backend 실행"
	@echo "make run-fe       - Frontend 실행"
	@echo "make stop-broker  - Mosquitto 브로커 중지"
	@echo "make clean        - 정리 (컨테이너, node_modules 등)"

# 프로젝트 초기 설정
setup:
	@echo "📦 Installing Backend dependencies..."
	cd backend && go mod download
	@echo "📦 Installing Frontend dependencies..."
	cd frontend && npm install
	@echo "✅ Setup complete!"

# Mosquitto 브로커 실행
run-broker:
	@echo "🚀 Starting Mosquitto broker..."
	docker-compose up -d
	@echo "✅ Mosquitto running on ports 1883 (MQTT), 9001 (WebSocket)"

# Backend 실행
run-be:
	@echo "🚀 Starting Backend..."
	cd backend && go run cmd/main.go

# Frontend 실행
run-fe:
	@echo "🚀 Starting Frontend..."
	cd frontend && npm start

# Mosquitto 브로커 중지
stop-broker:
	@echo "🛑 Stopping Mosquitto broker..."
	docker-compose down

# 정리
clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	rm -rf frontend/node_modules
	rm -rf backend/vendor
	@echo "✅ Cleanup complete!"
