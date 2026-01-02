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
        ├── cert-manager/       # cert-manager (TLS 인증서 관리)
        ├── nginx-gateway/      # NGINX Gateway Fabric
        └── gateway-routes/     # Gateway, HTTPRoute, TLS 리소스
            └── templates/
                ├── gateway.yaml
                ├── httproutes.yaml
                ├── certificate.yaml      # TLS 인증서 (선택적)
                └── clusterissuer.yaml    # Let's Encrypt Issuer (선택적)
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

## TLS/HTTPS 활성화 (선택적)

기본적으로 HTTP만 활성화되어 있습니다. HTTPS를 사용하려면 다음 단계를 따르세요:

### 1. values.yaml 수정

`charts/gateway/gateway-routes/values.yaml`:

```yaml
# TLS 활성화
tls:
  enabled: true  # false -> true로 변경

# HTTPS 리스너 주석 해제
gateway:
  listeners:
    - name: http
      # ...
    - name: https        # 주석 해제
      port: 443
      protocol: HTTPS
      # ...

# Let's Encrypt 설정 주석 해제
letsencrypt:
  email: your-email@example.com
  environment: staging  # 또는 prod

# Certificate 설정 주석 해제
certificate:
  name: echo-tls
  dnsNames:
    - echo.local
```

### 2. bootstrap 수정

`bootstrap/infra-gateway.yaml`에서 cert-manager 주석 해제:

```yaml
elements:
  - name: gateway-api-crds
  - name: cert-manager      # 주석 해제
    namespace: gateway
    path: cloud/ingress-gateway/charts/gateway/cert-manager
  - name: nginx-gateway
  - name: gateway-routes
```

### 3. 재배포

```bash
make deploy-gateway
```

## 사용된 컴포넌트

| 컴포넌트 | 버전 | 설명 |
|---------|------|------|
| Kind | v1.28.15 | 로컬 Kubernetes 클러스터 |
| ArgoCD | v7.8.28 | GitOps 기반 배포 도구 |
| echo-server | latest | 테스트용 Echo 서버 |
| NGINX Ingress Controller | v1.12.0 | Ingress 구현체 |
| Gateway API CRDs | v1.2.0 | Gateway API 리소스 정의 |
| NGINX Gateway Fabric | v2.2.1 | Gateway API 구현체 |
| cert-manager | v1.16.2 | TLS 인증서 자동 관리 (선택적) |

## 참고 자료

- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Gateway API](https://gateway-api.sigs.k8s.io/)
- [NGINX Gateway Fabric](https://docs.nginx.com/nginx-gateway-fabric/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [cert-manager](https://cert-manager.io/docs/)
