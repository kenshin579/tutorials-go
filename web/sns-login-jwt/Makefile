.DEFAULT_GOAL := help
.PHONY: help run-be run-fe run-all up down logs

# Google OAuth 자격 증명은 ~/.zshrc 의 WEB_SNS_LOGIN_GOOGLE_CLIENT_ID/SECRET 를 우선 사용하며,
# 없으면 backend/.env 의 GOOGLE_CLIENT_ID/SECRET 로 fallback 한다.

help: ## 사용 가능한 명령어 목록
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-8s\033[0m %s\n", $$1, $$2}'

run-be: ## Backend 실행 (Go, 포트 8080)
	cd backend && mkdir -p data && go run .

run-fe: ## Frontend 실행 (Vite, 포트 3000)
	cd frontend && npm run dev

run-all: ## Backend + Frontend 동시 실행 (Ctrl+C 로 함께 종료)
	@echo "▶ backend(:8080) + frontend(:3000) 실행 — Ctrl+C 로 종료"
	@trap 'kill 0' INT TERM EXIT; \
		( cd backend && mkdir -p data && go run . ) & \
		( cd frontend && npm run dev ) & \
		wait

up: ## docker compose 전체 스택 실행 (frontend :3000, backend :8080)
	docker compose up --build

down: ## docker compose 종료
	docker compose down

logs: ## docker compose 로그 확인
	docker compose logs -f
