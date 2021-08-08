REPOSITORY  := ingress-nginx
APP         := ingress-nginx
NAMESPACE   := echoserver

#https://kubernetes.github.io/ingress-nginx/deploy/

.PHONY: repo-add
repo-add:
	@helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	@helm repo update

.PHONY: helm-install
helm-install:
	@helm install $(APP) $(REPOSITORY)/$(APP) -n $(NAMESPACE)

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(APP) -n $(NAMESPACE)
