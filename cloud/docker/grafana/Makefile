GRAFANA_CONTAINER_NAME := go-grafana

.PHONY: grafana-create
grafana-create:
	@docker run --name $(GRAFANA_CONTAINER_NAME) -d \
	-e "GF_INSTALL_PLUGINS=grafana-clock-panel, grafana-simple-json-datasource" \
	-p 3000:3000 grafana/grafana-enterprise

.PHONY: grafana-delete
grafana-delete:
	@docker stop $(GRAFANA_CONTAINER_NAME)
	@docker rm $(GRAFANA_CONTAINER_NAME)

