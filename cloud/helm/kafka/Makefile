REPOSITORY  := bitnami
APP         := kafka
MYAPP         := kafka

.PHONY: helm-install
helm-install:
	helm3 install $(MYAPP) $(REPOSITORY)/$(APP) -f kafka-values.yaml

.PHONY: helm-uninstall
helm-uninstall:
	helm3 uninstall $(MYAPP)

