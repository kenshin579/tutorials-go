REPOSITORY  := ealenn
APP         := echo-server
NAMESPACE   := echoserver

.PHONY: repo-add
repo-add:
	@helm repo add ealenn https://ealenn.github.io/charts
	@helm repo update

.PHONY: helm-install
helm-install:
	helm upgrade -i $(APP) -f values.yaml $(REPOSITORY)/$(APP) -n $(NAMESPACE)

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(APP) -n $(NAMESPACE)

