# Go Delve Remote Debugging Example

Go Delve(dlv)를 사용한 원격 디버깅 샘플 프로젝트.
HTTP 서버 + 백그라운드 goroutine worker 구조로, 로컬/Docker/Kubernetes 환경에서 원격 디버깅을 실습할 수 있다.

## 프로젝트 구조

```
remote-debugging/
├── main.go                    # HTTP 서버 + goroutine worker
├── main_test.go               # 단위 테스트
├── Dockerfile.debug           # Delve 포함 디버그용 이미지
├── docker-compose.debug.yml   # Docker Compose (디버그)
├── k8s/
│   ├── deployment.yaml        # K8s Deployment
│   └── service.yaml           # K8s Service
└── README.md
```

## 로컬 실행

```bash
# 일반 실행
go run ./golang/debugging/remote-debugging/

# 테스트
go test -v ./golang/debugging/remote-debugging/
```

## 로컬 Delve 디버깅

```bash
# 소스코드에서 직접 디버깅
dlv debug ./golang/debugging/remote-debugging/

# 원격 디버깅 (headless 모드)
dlv debug ./golang/debugging/remote-debugging/ \
  --headless --listen=:2345 --api-version=2 --accept-multiclient

# 다른 터미널에서 연결
dlv connect localhost:2345
```

## Docker 디버깅

```bash
cd golang/debugging/remote-debugging/

# 빌드 및 실행
docker compose -f docker-compose.debug.yml up --build

# GoLand에서 Go Remote 설정으로 localhost:2345 연결
```

## Kubernetes 디버깅

```bash
# 이미지 빌드 (레포 루트에서)
docker build -f golang/debugging/remote-debugging/Dockerfile.debug -t go-debug-app:latest .

# 배포
kubectl apply -f golang/debugging/remote-debugging/k8s/

# 포트 포워딩
kubectl port-forward svc/go-debug-app 2345:2345 8080:8080

# GoLand에서 Go Remote 설정으로 localhost:2345 연결
```

## API 엔드포인트

| Method | Path | 설명 |
|--------|------|------|
| GET | `/health` | 헬스 체크 (worker 처리 건수 포함) |
| POST | `/process` | Job 생성 및 큐 등록 |

## 관련 블로그

- [Go Delve 원격 디버깅 완벽 가이드](https://blog-v2.advenoh.pe.kr)
