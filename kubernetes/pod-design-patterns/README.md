# Kubernetes Pod 디자인 패턴 실습

K8s Multi-Container Pod의 핵심 패턴들을 Kind 클러스터에서 직접 실습하는 예제 모음.

## 사전 준비물

- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

## 빠른 시작

```bash
# 클러스터 생성 + 이미지 빌드 + 로드
make all

# 패턴별 배포
make deploy-sidecar
make deploy-ambassador
make deploy-adapter
make deploy-init
make deploy-native

# 정리
make clean-pods   # Pod만 삭제
make clean        # 클러스터 삭제
```

## 프로젝트 구조

```
├── common/main-app/          # 공통 Go 웹 서버
├── sidecar/                  # Sidecar 패턴 (Go Request Logger)
├── ambassador/               # Ambassador 패턴 (Go Redis Proxy)
├── adapter/                  # Adapter 패턴 (Go Prometheus Exporter)
├── init-container/           # Init Container 패턴
└── native-sidecar/           # Native Sidecar (KEP-753)
```

## 패턴별 실습

### Sidecar - Request Logger

```bash
kubectl apply -f sidecar/sidecar-pod.yaml
kubectl wait --for=condition=Ready pod/sidecar-demo --timeout=60s

# 프록시를 통해 요청
kubectl exec sidecar-demo -c main-app -- wget -qO- http://localhost:8080/

# 로그 파일 확인
kubectl exec sidecar-demo -c main-app -- cat /var/log/app/access.log
```

### Ambassador - Redis Proxy

```bash
kubectl apply -f ambassador/redis-deployment.yaml
kubectl wait --for=condition=Available deployment/redis --timeout=60s
kubectl apply -f ambassador/ambassador-pod.yaml
kubectl wait --for=condition=Ready pod/ambassador-demo --timeout=60s

# 메인 앱에서 localhost:6379로 Redis 접근
kubectl exec ambassador-demo -c main-app -- sh -c 'echo PING | nc localhost 6379'
```

### Adapter - Prometheus Exporter

```bash
kubectl apply -f adapter/adapter-pod.yaml
kubectl wait --for=condition=Ready pod/adapter-demo --timeout=60s

# 메인 앱에 요청 몇 개 보내기
kubectl exec adapter-demo -c main-app -- wget -qO- http://localhost:3000/
kubectl exec adapter-demo -c main-app -- wget -qO- http://localhost:3000/

# Prometheus 형식 메트릭 확인
kubectl exec adapter-demo -c main-app -- wget -qO- http://localhost:9090/metrics
```

### Init Container 체이닝

```bash
# Redis 서비스가 먼저 필요 (Ambassador 예제와 공유)
kubectl apply -f ambassador/redis-deployment.yaml
kubectl wait --for=condition=Available deployment/redis --timeout=60s

kubectl apply -f init-container/init-chain-pod.yaml

# Pod 상태 변화 관찰
kubectl get pod init-chain-demo -w

# 설정 파일 확인
kubectl exec init-chain-demo -- cat /config/app-config.json
```

### Native Sidecar (K8s 1.28+)

```bash
kubectl apply -f native-sidecar/native-logger-pod.yaml
kubectl wait --for=condition=Ready pod/native-sidecar-demo --timeout=60s

# 요청 및 로그 확인
kubectl exec native-sidecar-demo -c main-app -- wget -qO- http://localhost:8080/
kubectl exec native-sidecar-demo -c main-app -- cat /var/log/app/access.log
```

## 관련 블로그

- [K8s Pod 디자인 패턴 (1) - Sidecar, Ambassador, Adapter](https://blog-v2.advenoh.pe.kr)
- [K8s Pod 디자인 패턴 (2) - Init Container 완벽 가이드](https://blog-v2.advenoh.pe.kr)
- [K8s Pod 디자인 패턴 (3) - Native Sidecar (KEP-753)](https://blog-v2.advenoh.pe.kr)
