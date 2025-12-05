.PHONY: all tf-init tf-apply tf-destroy clean

TERRAFORM_DIR := terraform

all: tf-init tf-apply

# Terraform commands
tf-init:
	cd $(TERRAFORM_DIR) && terraform init

tf-apply:
	cd $(TERRAFORM_DIR) && terraform apply -auto-approve

tf-destroy:
	cd $(TERRAFORM_DIR) && terraform destroy -auto-approve

tf-plan:
	cd $(TERRAFORM_DIR) && terraform plan

# Ingress deployment
deploy-ingress:
	kubectl apply -f bootstrap/apps.yaml
	kubectl apply -f bootstrap/infra-ingress.yaml

# Gateway deployment
deploy-gateway:
	kubectl apply -f bootstrap/apps.yaml
	kubectl apply -f bootstrap/infra-gateway.yaml

# Test commands
test-ingress:
	@echo "Testing Ingress..."
	@curl -s -H "Host: echo.local" http://localhost/ping | jq . || curl -s -H "Host: echo.local" http://localhost/ping

test-gateway:
	@echo "Testing Gateway..."
	@curl -s -H "Host: echo.local" http://localhost/ping | jq . || curl -s -H "Host: echo.local" http://localhost/ping

# ArgoCD access
argocd-password:
	@echo "ArgoCD Credentials:"
	@echo "  URL: http://localhost:8080"
	@echo "  Username: admin"
	@echo "  Password: password"

argocd-port-forward:
	kubectl port-forward svc/argocd-server -n argocd 8080:80

# Status check
status:
	@echo "=== Namespaces ==="
	kubectl get ns
	@echo "\n=== Pods in app namespace ==="
	kubectl get pods -n app
	@echo "\n=== Services ==="
	kubectl get svc -A | grep -E "(ingress|gateway|echo)"

# Cleanup
clean: tf-destroy

# Help
help:
	@echo "Available targets:"
	@echo "  tf-init        - Initialize Terraform"
	@echo "  tf-apply       - Create Kind cluster and install ArgoCD"
	@echo "  tf-destroy     - Destroy Kind cluster"
	@echo "  deploy-ingress - Deploy Ingress-based infrastructure"
	@echo "  deploy-gateway - Deploy Gateway API-based infrastructure"
	@echo "  test-ingress   - Test Ingress routing"
	@echo "  test-gateway   - Test Gateway routing"
	@echo "  argocd-password    - Show ArgoCD credentials"
	@echo "  argocd-port-forward - Port forward ArgoCD UI"
	@echo "  status         - Show cluster status"
	@echo "  clean          - Destroy all resources"
