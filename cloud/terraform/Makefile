CLUSTER_NAME := terraform-study-cluster

.PHONY: tf-init
tf-init:
	@echo "Initializing Terraform..."
	@terraform init -upgrade

.PHONY: tf-install
tf-install: tf-init
	@echo "Installing Terraform resources..."
	@terraform plan
	@terraform apply -auto-approve

.PHONY: tf-validate
tf-validate:
	@echo "Validating Terraform configuration..."
	@terraform validate

.PHONY: tf-destroy
tf-destroy:
	@echo "Destroying Terraform resources..."
	@terraform plan -destroy
	@terraform destroy -auto-approve

.PHONY: tf-clean
tf-clean:
	@echo "Cleaning up Terraform files..."
	@rm -rf .terraform/       # provider 캐시, module 다운로드 등
	@rm -f .terraform.*
	@rm -f terraform*

.PHONY: kind-delete
kind-delete: tf-clean
	@echo "Deleting Kind cluster..."
	@kind delete cluster --name $(CLUSTER_NAME)
