REPOSITORY  := bitnami
APP         := kafka
MYAPP       := kafka

.PHONY: helm-install
helm-install:
	helm install $(MYAPP) $(REPOSITORY)/$(APP) -f kafka-values.yaml

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(MYAPP)

