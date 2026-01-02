# Gateway vs Ingress 샘플 코드 분석 보고서

## 1. 개요

이 문서는 블로그 작성을 위해 만든 샘플 코드(`cloud/ingress-gateway`)와 실제 운영 중인 코드(`my-charts`)를 비교 분석한 결과입니다.

### 분석 대상
| 구분 | 경로 | 설명 |
|------|------|------|
| 샘플 코드 | `cloud/ingress-gateway/` | 블로그용 Gateway vs Ingress 예제 |
| 실제 코드 | `/Users/user/GolandProjects/my-charts/` | 운영 중인 Gateway 설정 |

---

## 2. 구조 비교

### 2.1 샘플 코드 구조
```
cloud/ingress-gateway/
├── bootstrap/
│   ├── infra-gateway.yaml      # Gateway 스택 ArgoCD ApplicationSet
│   ├── infra-ingress.yaml      # Ingress 스택 ArgoCD ApplicationSet
│   └── apps.yaml               # echo-server Application
├── charts/
│   ├── echo-server/            # 테스트용 백엔드 서비스
│   ├── gateway/
│   │   ├── gateway-api-crds/   # Gateway API CRD (placeholder)
│   │   ├── cert-manager/       # ✅ cert-manager 차트 (NEW)
│   │   ├── nginx-gateway/      # NGINX Gateway Fabric 설정
│   │   └── gateway-routes/     # Gateway + HTTPRoute + TLS 정의
│   │       └── templates/
│   │           ├── gateway.yaml
│   │           ├── httproutes.yaml
│   │           ├── certificate.yaml     # ✅ NEW
│   │           └── clusterissuer.yaml   # ✅ NEW
│   └── ingress/
│       ├── nginx-ingress/      # NGINX Ingress Controller 설정
│       └── ingress-routes/     # Ingress 리소스 정의
```

### 2.2 실제 코드 구조
```
my-charts/
├── bootstrap/
│   └── macmini-gateway.yaml    # Gateway 스택 ArgoCD ApplicationSet
├── charts/
│   ├── nginx-gateway/          # NGINX Gateway Fabric 설정
│   ├── cert-manager/           # 인증서 관리
│   └── gateway/                # Gateway 리소스 정의
│       └── templates/
│           ├── gateway.yaml
│           ├── httproutes.yaml
│           ├── certificate.yaml
│           ├── clusterissuer.yaml
│           ├── snippetsfilter.yaml
│           ├── basic-auth-secret.yaml
│           ├── backend-tls-policy.yaml
│           └── referencegrant.yaml
```

---

## 3. 상세 비교

### 3.1 Gateway 리소스 비교

#### 샘플 코드 (`gateway-routes/templates/gateway.yaml`) - ✅ 개선됨
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: {{ .Values.gateway.name }}
  namespace: {{ .Release.Namespace }}
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-{{ .Values.letsencrypt.environment }}
spec:
  gatewayClassName: {{ .Values.gateway.className }}
  listeners:
    {{- range .Values.gateway.listeners }}
    - name: {{ .name }}
      port: {{ .port }}
      protocol: {{ .protocol }}
      {{- if .tls }}
      tls:
        mode: {{ .tls.mode }}
        certificateRefs:
          {{- range .tls.certificateRefs }}
          - kind: {{ .kind }}
            name: {{ .name }}
          {{- end }}
      {{- end }}
      allowedRoutes:
        namespaces:
          from: {{ .allowedRoutes.from }}
    {{- end }}
```

#### 실제 코드 (`gateway/templates/gateway.yaml`)
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: {{ .Values.gateway.name }}
  namespace: {{ .Release.Namespace }}
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-{{ .Values.letsencrypt.environment }}
spec:
  gatewayClassName: {{ .Values.gateway.gatewayClassName }}
  listeners:
  {{- range .Values.gateway.listeners }}
  - name: {{ .name }}
    protocol: {{ .protocol }}
    port: {{ .port }}
    {{- if .tls }}
    tls:
      mode: {{ .tls.mode }}
      certificateRefs:
      {{- range .tls.certificateRefs }}
      - kind: {{ .kind }}
        name: {{ .name }}
      {{- end }}
    {{- end }}
  {{- end }}
```

#### 🔍 차이점 분석

| 항목 | 샘플 코드 | 실제 코드 | 상태 |
|------|-----------|-----------|------|
| TLS 설정 | ✅ TLS 지원 | ✅ TLS 지원 | ✅ 동일 |
| cert-manager 연동 | ✅ annotation 연동 | ✅ annotation 연동 | ✅ 동일 |
| allowedRoutes | ✅ 있음 | ❌ 없음 (기본값 사용) | ✅ 샘플이 더 명시적 |

---

### 3.2 HTTPRoute 비교

#### 샘플 코드 (`gateway-routes/templates/httproutes.yaml`)
```yaml
{{- range .Values.httpRoutes }}
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: {{ .name }}
spec:
  parentRefs:
    - name: {{ $.Values.gateway.name }}
  hostnames:
    {{- range .hostnames }}
    - {{ . }}
    {{- end }}
  rules:
    {{- range .rules }}
    - matches:
        {{- range .matches }}
        - path:
            type: {{ .path.type }}
            value: {{ .path.value }}
        {{- end }}
      backendRefs:
        {{- range .backendRefs }}
        - name: {{ .name }}
          namespace: {{ .namespace }}
          port: {{ .port }}
        {{- end }}
    {{- end }}
{{- end }}
```

#### 실제 코드 (`gateway/templates/httproutes.yaml`)
```yaml
{{- range .Values.routes }}
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: {{ .name }}
spec:
  parentRefs:
  - name: {{ $.Values.gateway.name }}
  hostnames:
  - {{ .hostname }}
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: {{ .path | default "/" }}
    {{- if .basicAuth }}
    filters:
    - type: ExtensionRef
      extensionRef:
        group: gateway.nginx.org
        kind: SnippetsFilter
        name: {{ .name }}-basic-auth
    {{- end }}
    backendRefs:
    - name: {{ .service.name }}
      namespace: {{ .service.namespace }}
      port: {{ .service.port }}
{{- end }}
```

#### 🔍 차이점 분석

| 항목 | 샘플 코드 | 실제 코드 | 상태 |
|------|-----------|-----------|------|
| hostname 배열 | ✅ 다중 hostname | 단일 hostname | ✅ 샘플이 더 유연 |
| 다중 matches | ✅ range로 반복 | 단일 match | ✅ 샘플이 더 유연 |
| 다중 backendRefs | ✅ range로 반복 | 단일 backendRef | ✅ 샘플이 더 유연 |
| Basic Auth | ❌ 없음 | ✅ SnippetsFilter 연동 | 선택적 기능 |
| filters | ❌ 없음 | ✅ 지원 | 선택적 기능 |

---

### 3.3 cert-manager 설정 비교 (NEW)

#### 샘플 코드 (`cert-manager/values.yaml`) - ✅ 추가됨
```yaml
cert-manager:
  crds:
    enabled: true
    keep: true

  # Gateway API 지원 활성화
  featureGates: "ExperimentalGatewayAPISupport=true"

  config:
    apiVersion: controller.config.cert-manager.io/v1alpha1
    kind: ControllerConfiguration
    enableGatewayAPI: true

  resources:
    requests:
      cpu: 10m
      memory: 32Mi
    limits:
      cpu: 100m
      memory: 128Mi
```

#### 실제 코드 (`cert-manager/values.yaml`)
```yaml
cert-manager:
  crds:
    enabled: true
    keep: true

  featureGates: "ExperimentalGatewayAPISupport=true"

  config:
    apiVersion: controller.config.cert-manager.io/v1alpha1
    kind: ControllerConfiguration
    enableGatewayAPI: true
  # ... (동일한 설정)
```

| 항목 | 샘플 코드 | 실제 코드 | 상태 |
|------|-----------|-----------|------|
| Gateway API 지원 | ✅ | ✅ | ✅ 동일 |
| CRD 설치 | ✅ | ✅ | ✅ 동일 |
| 리소스 제한 | ✅ | ✅ | ✅ 동일 |

---

### 3.4 Certificate & ClusterIssuer 비교 (NEW)

#### 샘플 코드 (`gateway-routes/templates/clusterissuer.yaml`) - ✅ 추가됨
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-{{ .Values.letsencrypt.environment }}
spec:
  acme:
    email: {{ .Values.letsencrypt.email }}
    server: {{ if eq .Values.letsencrypt.environment "prod" }}https://acme-v02.api.letsencrypt.org/directory{{ else }}https://acme-staging-v02.api.letsencrypt.org/directory{{ end }}
    privateKeySecretRef:
      name: letsencrypt-{{ .Values.letsencrypt.environment }}-key
    solvers:
      - http01:
          gatewayHTTPRoute:
            parentRefs:
              - name: {{ .Values.gateway.name }}
                namespace: {{ .Release.Namespace }}
                kind: Gateway
```

#### 샘플 코드 (`gateway-routes/templates/certificate.yaml`) - ✅ 추가됨
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: {{ .Values.certificate.name }}
  namespace: {{ .Release.Namespace }}
spec:
  secretName: {{ .Values.certificate.name }}
  dnsNames:
  {{- range .Values.certificate.dnsNames }}
  - {{ . }}
  {{- end }}
  issuerRef:
    name: letsencrypt-{{ .Values.letsencrypt.environment }}
    kind: ClusterIssuer
```

| 항목 | 샘플 코드 | 실제 코드 | 상태 |
|------|-----------|-----------|------|
| ClusterIssuer | ✅ | ✅ | ✅ 동일 |
| Certificate | ✅ | ✅ | ✅ 동일 |
| HTTP-01 챌린지 | ✅ Gateway 연동 | ✅ Gateway 연동 | ✅ 동일 |

---

### 3.5 NGINX Gateway Fabric 설정 비교

#### 샘플 코드 (`nginx-gateway/values.yaml`)
```yaml
nginx-gateway-fabric:
  service:
    type: NodePort
  nginxGateway:
    gwAPIExperimentalFeatures:
      enable: true
  nodeSelector:
    ingress-ready: "true"
  tolerations:
    - key: "node-role.kubernetes.io/control-plane"
      operator: "Equal"
      effect: "NoSchedule"
```

#### 실제 코드 (`nginx-gateway/values.yaml`)
```yaml
nginx-gateway-fabric:
  nginxGateway:
    productTelemetry:
      enable: false
    gatewayClassName: nginx
    gwAPIExperimentalFeatures:
      enable: true
    snippetsFilters:
      enable: true                    # ✅ SnippetsFilter 활성화
    nodeSelector:
      node-role.kubernetes.io/control-plane: ""
    tolerations:
      - key: node-role.kubernetes.io/control-plane
        operator: Exists
        effect: NoSchedule

  nginx:
    pod:
      volumes:
        - name: htpasswd
          secret:
            secretName: basic-auth-htpasswd   # ✅ Basic Auth 시크릿 마운트
    container:
      hostPorts:
        - port: 80
          containerPort: 80
        - port: 443
          containerPort: 443
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
      volumeMounts:
        - name: htpasswd
          mountPath: /etc/nginx/auth
          readOnly: true

  service:
    type: NodePort
    ports:
      - port: 80
        nodePort: 30026
      - port: 443
        nodePort: 30027
```

#### 🔍 차이점 분석

| 항목 | 샘플 코드 | 실제 코드 | 상태 |
|------|-----------|-----------|------|
| SnippetsFilter | ❌ 비활성화 | ✅ 활성화 | 선택적 (Basic Auth 필요시) |
| hostPorts | ❌ 없음 | ✅ 80/443 포트 | 선택적 (Kind 환경) |
| htpasswd 볼륨 | ❌ 없음 | ✅ 마운트 설정 | 선택적 (Basic Auth 필요시) |
| 리소스 제한 | ❌ 없음 | ✅ requests/limits | 권장 |
| nodePort 지정 | ❌ 자동 할당 | ✅ 명시적 지정 | 선택적 |

---

### 3.6 샘플 코드에 없는 기능 (실제 코드에만 존재)

| 기능 | 파일 | 설명 | 필요성 |
|------|------|------|--------|
| **Basic Auth** | `snippetsfilter.yaml` | NGINX SnippetsFilter를 통한 인증 | 선택적 |
| **Basic Auth Secret** | `basic-auth-secret.yaml` | htpasswd 시크릿 | 선택적 |
| **Backend TLS Policy** | `backend-tls-policy.yaml` | 백엔드 서비스 TLS 설정 | 선택적 |
| **ReferenceGrant** | `referencegrant.yaml` | 크로스 네임스페이스 참조 허용 | 선택적 |

---

## 4. 결론 및 권장사항

### 4.1 샘플 코드 평가

#### ✅ 완료된 개선 사항
1. **TLS 설정 추가**: Gateway에 HTTPS 리스너 추가 완료
2. **cert-manager 연동**: ClusterIssuer, Certificate 템플릿 추가 완료
3. **cert-manager 차트**: Gateway API 지원 활성화된 차트 추가 완료
4. **bootstrap 업데이트**: cert-manager 배포 순서 포함

#### ✅ 기존 장점
1. **Ingress vs Gateway 구조 분리**: 명확하게 두 방식을 비교할 수 있는 구조
2. **유연한 HTTPRoute 템플릿**: 다중 hostname, matches, backendRefs 지원
3. **Kind 클러스터 호환**: nodeSelector, tolerations 설정 적절
4. **ArgoCD ApplicationSet 활용**: 배포 자동화 구조 잘 설계됨

#### 📋 선택적 개선 사항 (필요시 추가)

| 우선순위 | 항목 | 설명 |
|----------|------|------|
| 🟡 선택적 | **SnippetsFilter 예제** | Basic Auth 같은 실용적 예제 |
| 🟡 선택적 | **hostPort 설정** | Kind 클러스터에서 직접 접근 설정 |
| 🟢 선택적 | **리소스 제한** | NGINX Gateway Pod 리소스 설정 |

### 4.2 블로그 작성 시 참고사항

1. **HTTP + HTTPS 예제 모두 포함**: 현재 샘플로 두 가지 모두 설명 가능
2. **Let's Encrypt 연동**: staging/prod 환경 모두 설명 가능
3. **실제 동작 확인 필요**: echo-server로 라우팅 테스트 후 스크린샷 추가 권장

---

## 5. 적용된 개선 사항

### 5.1 Gateway TLS 리스너 (✅ 적용됨)

`gateway-routes/values.yaml`:
```yaml
gateway:
  name: echo-gateway
  className: nginx
  listeners:
    - name: http
      port: 80
      protocol: HTTP
      allowedRoutes:
        from: All
    - name: https
      port: 443
      protocol: HTTPS
      allowedRoutes:
        from: All
      tls:
        mode: Terminate
        certificateRefs:
          - kind: Secret
            name: echo-tls

letsencrypt:
  email: your-email@example.com
  environment: staging

certificate:
  name: echo-tls
  dnsNames:
    - echo.local
```

### 5.2 cert-manager 차트 (✅ 적용됨)

`gateway/cert-manager/Chart.yaml`:
```yaml
apiVersion: v2
name: cert-manager
description: cert-manager for automatic TLS certificate management
type: application
version: 1.0.0
appVersion: "v1.16.2"
dependencies:
  - name: cert-manager
    version: v1.16.2
    repository: https://charts.jetstack.io
```

### 5.3 ArgoCD 배포 순서 (✅ 적용됨)

`bootstrap/infra-gateway.yaml`:
```yaml
elements:
  - name: gateway-api-crds    # 1. CRD 먼저 설치
  - name: cert-manager        # 2. cert-manager 설치
  - name: nginx-gateway       # 3. NGINX Gateway Fabric 설치
  - name: gateway-routes      # 4. Gateway + HTTPRoute + Certificate 설치
```

---

## 6. 요약

| 구분 | 샘플 코드 (개선 후) | 실제 코드 |
|------|---------------------|-----------|
| **목적** | 교육/블로그 예제 | 프로덕션 운영 |
| **TLS** | ✅ Let's Encrypt | ✅ Let's Encrypt |
| **cert-manager** | ✅ Gateway API 연동 | ✅ Gateway API 연동 |
| **인증** | ❌ 없음 (선택적) | ✅ Basic Auth |
| **구조** | 단순하고 명확 | 복잡하지만 완전함 |
| **사용성** | 학습에 적합 | 운영에 적합 |

**최종 평가**: 샘플 코드가 TLS와 cert-manager 연동이 추가되어 프로덕션 수준의 HTTPS 예제를 보여줄 수 있게 되었습니다. Gateway vs Ingress 개념 비교와 함께 실제 HTTPS 구현까지 블로그에서 다룰 수 있습니다.
