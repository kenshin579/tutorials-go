CLUSTER_NAME := pod-patterns
IMAGES := main-app request-logger redis-proxy metrics-adapter

.PHONY: cluster-up cluster-down build load all clean

cluster-up:
	kind create cluster --name $(CLUSTER_NAME) --config kind-config.yaml

cluster-down:
	kind delete cluster --name $(CLUSTER_NAME)

build:
	docker build -t main-app:local common/main-app/
	docker build -t request-logger:local sidecar/request-logger/
	docker build -t redis-proxy:local ambassador/redis-proxy/
	docker build -t metrics-adapter:local adapter/metrics-adapter/

load: build
	@for img in $(IMAGES); do \
		kind load docker-image $$img:local --name $(CLUSTER_NAME); \
	done

all: cluster-up load

clean: cluster-down

# 개별 패턴 배포
deploy-sidecar:
	kubectl apply -f sidecar/sidecar-pod.yaml

deploy-ambassador:
	kubectl apply -f ambassador/redis-deployment.yaml
	kubectl apply -f ambassador/ambassador-pod.yaml

deploy-adapter:
	kubectl apply -f adapter/adapter-pod.yaml

deploy-init:
	kubectl apply -f init-container/init-chain-pod.yaml

deploy-native:
	kubectl apply -f native-sidecar/native-logger-pod.yaml

# 전체 정리
clean-pods:
	kubectl delete pod --all --grace-period=0 --force 2>/dev/null || true
	kubectl delete deployment --all 2>/dev/null || true
	kubectl delete service redis-service 2>/dev/null || true
	kubectl delete job --all 2>/dev/null || true
