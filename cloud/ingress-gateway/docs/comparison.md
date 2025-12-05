# Ingress vs Gateway API 비교

## 개요

Kubernetes에서 외부 트래픽을 클러스터 내부 서비스로 라우팅하는 두 가지 주요 방식을 비교합니다.

## 리소스 구조 비교

### Ingress 방식

```yaml
# 단일 Ingress 리소스로 모든 설정
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

### Gateway API 방식

```yaml
# 1. Gateway 리소스 (인프라 팀 관리)
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
---
# 2. HTTPRoute 리소스 (개발 팀 관리)
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: echo-server-route
spec:
  parentRefs:
    - name: echo-gateway
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

## 기능 비교표

| 기능 | Ingress | Gateway API |
|-----|---------|-------------|
| **기본 HTTP 라우팅** | O | O |
| **경로 기반 라우팅** | O | O |
| **호스트 기반 라우팅** | O | O |
| **TLS 종료** | O | O |
| **헤더 기반 라우팅** | Annotation 의존 | 표준 지원 |
| **트래픽 분할** | Annotation 의존 | 표준 지원 |
| **요청/응답 수정** | Annotation 의존 | 표준 지원 |
| **TCP/UDP 라우팅** | X | O |
| **gRPC 라우팅** | Annotation 의존 | 표준 지원 |
| **역할 기반 분리** | X | O (Gateway/Route) |

## 장단점 비교

### Ingress

**장점:**
- 간단한 설정 - 단일 리소스로 모든 설정 가능
- 널리 사용됨 - 풍부한 문서와 커뮤니티
- 학습 곡선 낮음 - 빠르게 시작 가능

**단점:**
- Annotation 의존성 - 구현체마다 다른 설정 방식
- 제한된 프로토콜 - HTTP/HTTPS만 지원
- 표준화 부족 - 구현체별 동작 차이

### Gateway API

**장점:**
- 표준화된 API - 구현체 간 일관된 동작
- 역할 기반 분리 - 인프라/앱 팀 분리 가능
- 다양한 프로토콜 - TCP, UDP, gRPC 등 지원
- 확장 가능한 구조 - CRD 기반 확장

**단점:**
- 상대적으로 새로움 - 일부 기능 실험적
- 설정 복잡성 - 여러 리소스 관리 필요
- 구현체 성숙도 - 일부 구현체 미성숙

## 사용 시나리오

### Ingress 추천 상황
- 간단한 HTTP/HTTPS 라우팅만 필요한 경우
- 빠르게 프로토타입을 만들어야 하는 경우
- 기존 Ingress 인프라가 있는 경우

### Gateway API 추천 상황
- 다양한 프로토콜(TCP, UDP, gRPC) 지원이 필요한 경우
- 인프라 팀과 개발 팀의 역할 분리가 필요한 경우
- 복잡한 트래픽 관리(분할, 미러링 등)가 필요한 경우
- 새로운 프로젝트를 시작하는 경우

## 마이그레이션 고려사항

### Ingress → Gateway API

1. **점진적 마이그레이션**: 두 방식을 동시에 운영 가능
2. **호환성 레이어**: 일부 구현체는 Ingress를 Gateway API로 변환
3. **설정 매핑**: Annotation → 표준 필드 매핑 필요

## 결론

- **Ingress**: 단순함이 필요한 경우 여전히 유효한 선택
- **Gateway API**: 미래 지향적이고 확장 가능한 표준, 새 프로젝트에 권장
