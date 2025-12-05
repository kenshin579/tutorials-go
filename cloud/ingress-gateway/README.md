# Ingress vs Gateway API Demo

Kubernetes에서 외부 트래픽을 처리하는 두 가지 방식인 **Ingress**와 **Gateway API**를 비교하는 샘플 프로젝트입니다.

## 아키텍처

```
                    ┌─────────────────────────────────────────┐
                    │           Kind Cluster                  │
                    │  ┌─────────────────────────────────┐   │
                    │  │           ArgoCD                 │   │
                    │  └─────────────────────────────────┘   │
                    │                                         │
  ┌──────────┐      │  ┌─────────────────────────────────┐   │
  │  Client  │──────│──│  Ingress or Gateway API         │   │
  └──────────┘      │  └─────────────────────────────────┘   │
                    │                 │                       │
                    │                 ▼                       │
                    │  ┌─────────────────────────────────┐   │
                    │  │         echo-server              │   │
                    │  │    (kenshin579/echo-server:latest) │   │
                    │  └─────────────────────────────────┘   │
                    └─────────────────────────────────────────┘
```

## 사전 요구사항

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Docker](https://www.docker.com/get-started)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

## 빠른 시작

### 1. 클러스터 생성

```bash
make tf-init
make tf-apply
```

### 2-A. Ingress 방식으로 배포

```bash
make deploy-ingress
```

### 2-B. Gateway API 방식으로 배포

```bash
make deploy-gateway
```

### 3. 테스트

```bash
# Ingress 테스트
make test-ingress

# Gateway 테스트
make test-gateway

# 직접 curl 호출
curl -H "Host: echo.local" http://localhost/ping
```

### 4. ArgoCD 접속

```bash
# 포트 포워딩
make argocd-port-forward

# 접속 정보 확인
make argocd-password
```

- URL: http://localhost:8080
- Username: `admin`
- Password: `password`

### 5. 정리

```bash
make clean
```

## 디렉토리 구조

```
.
├── Makefile                    # 자동화 명령어
├── README.md
├── terraform/                  # Terraform 인프라 코드
│   ├── main.tf                 # Provider 설정
│   ├── variables.tf            # 변수 정의
│   ├── outputs.tf              # 출력 정의
│   ├── k8s.tf                  # Kind 클러스터 설정
│   └── modules/infra/          # ArgoCD 설치 모듈
├── bootstrap/                  # ArgoCD Applications
│   ├── apps.yaml               # echo-server 배포
│   ├── infra-ingress.yaml      # Ingress 인프라
│   └── infra-gateway.yaml      # Gateway 인프라
└── charts/                     # Helm Charts
    ├── echo-server/            # 샘플 애플리케이션
    ├── ingress/                # Ingress 관련 차트
    │   ├── nginx-ingress/      # NGINX Ingress Controller
    │   └── ingress-routes/     # Ingress 리소스
    └── gateway/                # Gateway API 관련 차트
        ├── gateway-api-crds/   # Gateway API CRDs
        ├── nginx-gateway/      # NGINX Gateway Fabric
        └── gateway-routes/     # Gateway, HTTPRoute 리소스
```

## 비교: Ingress vs Gateway API

| 항목 | Ingress | Gateway API |
|------|---------|-------------|
| API 버전 | networking.k8s.io/v1 | gateway.networking.k8s.io/v1 |
| 리소스 수 | 1개 (Ingress) | 2개+ (Gateway, HTTPRoute) |
| 역할 분리 | 없음 (단일 리소스) | 명확함 (인프라/앱 분리) |
| 확장성 | Annotation 기반 | 표준화된 CRD |
| 프로토콜 지원 | HTTP/HTTPS | HTTP, HTTPS, TCP, UDP, gRPC |
| 표준화 | 구현체마다 상이 | Kubernetes SIG 표준 |

## 사용된 컴포넌트

| 컴포넌트 | 버전 | 설명 |
|---------|------|------|
| Kind | v1.28.15 | 로컬 Kubernetes 클러스터 |
| ArgoCD | v7.8.28 | GitOps 기반 배포 도구 |
| echo-server | latest | 테스트용 Echo 서버 |
| NGINX Ingress Controller | v1.12.0 | Ingress 구현체 |
| Gateway API CRDs | v1.2.0 | Gateway API 리소스 정의 |
| NGINX Gateway Fabric | v2.2.1 | Gateway API 구현체 |

## 참고 자료

- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Gateway API](https://gateway-api.sigs.k8s.io/)
- [NGINX Gateway Fabric](https://docs.nginx.com/nginx-gateway-fabric/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
