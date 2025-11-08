.PHONY: help init plan apply destroy clean test logs status

help: ## 도움말 표시
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## Terraform 초기화
	terraform init

plan: ## 배포 계획 확인
	terraform plan

apply: ## Ollama 클러스터 배포
	terraform apply -auto-approve

destroy: ## 클러스터 삭제
	terraform destroy -auto-approve

clean: destroy ## 모든 리소스 정리
	rm -rf .terraform*
	rm -f terraform.tfstate*

test: ## Ollama API 테스트
	@echo "Ollama API 상태 확인..."
	@curl -s http://localhost:30025/api/tags | jq '.' || echo "API가 아직 준비되지 않았습니다."

logs: ## Ollama Pod 로그 확인
	kubectl logs -n ollama -l app.kubernetes.io/name=ollama -f

status: ## 클러스터 상태 확인
	@echo "=== Nodes ==="
	@kubectl get nodes
	@echo "\n=== Pods ==="
	@kubectl get pods -n ollama
	@echo "\n=== Services ==="
	@kubectl get svc -n ollama

model-pull: ## 모델 다운로드 (qwen2.5:0.5b)
	curl -X POST http://localhost:30025/api/pull -d '{"name": "qwen2.5:0.5b"}'

model-list: ## 설치된 모델 목록
	@curl -s http://localhost:30025/api/tags | jq '.models[] | {name: .name, size: .size}'

chat: ## 간단한 채팅 테스트
	@echo "테스트 메시지 전송 중..."
	@curl -s -X POST http://localhost:30025/api/generate \
		-d '{"model": "qwen2.5:0.5b", "prompt": "Hello! How are you?", "stream": false}' \
		| jq -r '.response' 