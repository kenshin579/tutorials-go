# Ingress vs Gateway API 샘플 프로젝트 PRD

## 1. 프로젝트 개요

### 목적
Kubernetes에서 외부 트래픽을 처리하는 두 가지 방식인 **Ingress**와 **Gateway API**를 비교하는 샘플 코드를 작성하여 블로그 자료로 활용한다.

### 참고 프로젝트
- **my-charts**: `/Users/user/GolandProjects/my-charts`
  - Terraform + ArgoCD + Gateway API 기반 프로덕션 인프라 구성
  - Kind 클러스터, NGINX Gateway Fabric, cert-manager 사용
- **echo-server**: `https://github.com/kenshin579/echo-server`
  - 샘플 애플리케이션으로 사용할 Go 기반 Echo Server
  - Docker 이미지: `kenshin579/echo-server:latest`

### 작성 위치
- `cloud/ingress-gateway/`

---

## 2. 디렉토리 구조

```
cloud/ingress-gateway/
├── README.md                          # 프로젝트 설명 및 실행 가이드
├── Makefile                           # 자동화 명령어
│
├── terraform/                         # Terraform 인프라 코드
│   ├── main.tf                        # Provider 설정
│   ├── variables.tf                   # 변수 정의
│   ├── outputs.tf                     # 출력 정의
│   ├── k8s.tf                         # Kind 클러스터 설정
│   └── modules/
│       └── infra/
│           └── infra.tf               # ArgoCD 설치
│
├── bootstrap/                         # ArgoCD ApplicationSet
│   ├── apps.yaml                      # 샘플 앱 배포
│   ├── infra-ingress.yaml            # Ingress 기반 인프라
│   └── infra-gateway.yaml            # Gateway API 기반 인프라
│
├── charts/                            # Helm Charts
│   │
│   ├── echo-server/                   # 샘플 애플리케이션 (echo-server)
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       ├── deployment.yaml
│   │       ├── service.yaml
│   │       └── _helpers.tpl
│   │
│   ├── ingress/                       # Ingress 방식
│   │   ├── nginx-ingress/             # NGINX Ingress Controller
│   │   │   ├── Chart.yaml
│   │   │   └── values.yaml
│   │   └── ingress-routes/            # Ingress 리소스 정의
│   │       ├── Chart.yaml
│   │       ├── values.yaml
│   │       └── templates/
│   │           └── ingress.yaml
│   │
│   └── gateway/                       # Gateway API 방식
│       ├── gateway-api-crds/          # Gateway API CRDs
│       │   ├── Chart.yaml
│       │   └── values.yaml
│       ├── nginx-gateway/             # NGINX Gateway Fabric
│       │   ├── Chart.yaml
│       │   └── values.yaml
│       └── gateway-routes/            # Gateway, HTTPRoute 정의
│           ├── Chart.yaml
│           ├── values.yaml
│           └── templates/
│               ├── gateway.yaml
│               └── httproutes.yaml
│
└── docs/                              # 추가 문서
    └── comparison.md                  # Ingress vs Gateway 비교 표
```

---

## 3. 핵심 컴포넌트

### 3.1 공통 인프라

| 컴포넌트 | 설명 | 버전/이미지 |
|---------|------|------------|
| Kind Cluster | 로컬 Kubernetes 클러스터 | kindest/node:v1.28.15 |
| ArgoCD | GitOps 기반 배포 도구 | v7.8.28 |
| Sample App | 테스트용 echo-server | kenshin579/echo-server:latest |

### 3.2 Ingress 방식

| 컴포넌트 | 설명 | 버전 |
|---------|------|------|
| NGINX Ingress Controller | Ingress 구현체 | v1.12.0 |
| Ingress 리소스 | 라우팅 규칙 정의 | networking.k8s.io/v1 |

### 3.3 Gateway API 방식

| 컴포넌트 | 설명 | 버전 |
|---------|------|------|
| Gateway API CRDs | Gateway API 리소스 정의 | v1.2.0 |
| NGINX Gateway Fabric | Gateway API 구현체 | v2.2.1 |
| Gateway | 게이트웨이 리소스 | gateway.networking.k8s.io/v1 |
| HTTPRoute | HTTP 라우팅 규칙 | gateway.networking.k8s.io/v1 |

---

## 4. 구현 요구사항

### 4.1 Terraform 구성

#### Kind 클러스터 설정
- **노드 구성**: control-plane 1개 + worker 1개 (간소화)
- **포트 매핑**:
  - 80: HTTP 트래픽
  - 443: HTTPS 트래픽 (선택)
  - 30080: NodePort (대체 접근)

#### Provider 설정
```hcl
required_providers {
  kind = {
    source  = "tehcyx/kind"
    version = "0.8.0"
  }
  kubernetes = {
    source  = "hashicorp/kubernetes"
    version = "~> 2.35"
  }
  helm = {
    source  = "hashicorp/helm"
    version = "~> 2.17"
  }
}
```

### 4.2 ArgoCD 설정

#### ApplicationSet 구성
- **앱 배포**: echo-server
- **인프라 배포**: ingress 또는 gateway 선택적 배포
- **Auto-sync**: 활성화 (prune, selfHeal)

### 4.3 샘플 애플리케이션

#### Echo Server 설정
- **소스**: https://github.com/kenshin579/echo-server
- **이미지**: `kenshin579/echo-server:latest`
- **포트**: 80
- **Replicas**: 1
- **네임스페이스**: `app`
- **서비스**: ClusterIP
- **엔드포인트**:
  - `/ping`: Health check
  - `/swagger/index.html`: API 문서

### 4.4 Ingress 구성

#### NGINX Ingress Controller
- Helm Chart: `ingress-nginx/ingress-nginx`
- **서비스 타입**: NodePort (Kind 호환)
- **네임스페이스**: `ingress-nginx`

#### Ingress 리소스
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: echo-server-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: echo.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: echo-server
                port:
                  number: 80
```

### 4.5 Gateway API 구성

#### Gateway API CRDs
- Source: `kubernetes-sigs/gateway-api`
- 버전: v1.2.0
- Experimental CRDs 포함

#### NGINX Gateway Fabric
- Helm Chart: `oci://ghcr.io/nginx/charts/nginx-gateway-fabric`
- **서비스 타입**: NodePort (Kind 호환)
- **네임스페이스**: `gateway`

#### Gateway 리소스
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: echo-gateway
spec:
  gatewayClassName: nginx
  listeners:
    - name: http
      port: 80
      protocol: HTTP
      allowedRoutes:
        namespaces:
          from: All
```

#### HTTPRoute 리소스
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: echo-server-route
spec:
  parentRefs:
    - name: echo-gateway
      namespace: gateway
  hostnames:
    - echo.local
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: echo-server
          port: 80
```

---

## 5. 보안 설정

### 인증 정보 (민감 정보 대체값)
| 항목 | 값 |
|-----|-----|
| Admin Username | admin |
| Admin Password | password |
| ArgoCD Password | password |

---

## 6. 실행 시나리오

### 6.1 Ingress 방식 테스트
```bash
# 1. 클러스터 생성 + ArgoCD 설치
make tf-init
make tf-apply

# 2. Ingress 인프라 배포
kubectl apply -f bootstrap/infra-ingress.yaml
kubectl apply -f bootstrap/apps.yaml

# 3. 테스트
curl -H "Host: echo.local" http://localhost:80

# 4. 정리
make tf-destroy
```

### 6.2 Gateway API 방식 테스트
```bash
# 1. 클러스터 생성 + ArgoCD 설치
make tf-init
make tf-apply

# 2. Gateway 인프라 배포
kubectl apply -f bootstrap/infra-gateway.yaml
kubectl apply -f bootstrap/apps.yaml

# 3. 테스트
curl -H "Host: echo.local" http://localhost:80

# 4. 정리
make tf-destroy
```

---

## 7. 비교 포인트 (블로그용)

### 아키텍처 비교

| 항목 | Ingress | Gateway API |
|------|---------|-------------|
| API 버전 | networking.k8s.io/v1 | gateway.networking.k8s.io/v1 |
| 리소스 수 | 1개 (Ingress) | 2개+ (Gateway, HTTPRoute) |
| 역할 분리 | 없음 (단일 리소스) | 명확함 (인프라/앱 분리) |
| 확장성 | Annotation 기반 | 표준화된 CRD |
| 프로토콜 지원 | HTTP/HTTPS | HTTP, HTTPS, TCP, UDP, gRPC |
| 표준화 | 구현체마다 상이 | Kubernetes SIG 표준 |

### 장단점 비교

#### Ingress
**장점**:
- 간단한 설정
- 널리 사용됨
- 학습 곡선 낮음

**단점**:
- Annotation 의존성
- 구현체별 차이
- 제한된 프로토콜

#### Gateway API
**장점**:
- 표준화된 API
- 역할 기반 분리
- 다양한 프로토콜
- 확장 가능한 구조

**단점**:
- 상대적으로 새로움
- 설정 복잡성
- 일부 구현체 미성숙

---

## 8. 참고 자료

### my-charts 프로젝트 주요 파일
- `k8s.tf`: Kind 클러스터 구성 참고
- `modules/infra/infra.tf`: ArgoCD 설치 참고
- `bootstrap/macmini-gateway.yaml`: Gateway ApplicationSet 참고
- `charts/gateway/`: Gateway 리소스 구성 참고
- `charts/nginx-gateway/`: NGINX Gateway Fabric 구성 참고

### echo-server 프로젝트
- GitHub: https://github.com/kenshin579/echo-server
- Docker Hub: `kenshin579/echo-server:latest`
- `deploy/echo-deployments.yaml`: Kubernetes 배포 예시 참고

### 공식 문서
- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Gateway API](https://gateway-api.sigs.k8s.io/)
- [NGINX Gateway Fabric](https://docs.nginx.com/nginx-gateway-fabric/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
