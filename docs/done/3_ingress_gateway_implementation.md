# Ingress vs Gateway API 구현 문서

## 1. 프로젝트 구조

```
cloud/ingress-gateway/
├── README.md
├── Makefile
├── terraform/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── k8s.tf
│   └── modules/infra/infra.tf
├── bootstrap/
│   ├── apps.yaml
│   ├── infra-ingress.yaml
│   └── infra-gateway.yaml
└── charts/
    ├── echo-server/
    ├── ingress/
    │   ├── nginx-ingress/
    │   └── ingress-routes/
    └── gateway/
        ├── gateway-api-crds/
        ├── nginx-gateway/
        └── gateway-routes/
```

---

## 2. Terraform 구현

### 2.1 main.tf - Provider 설정

```hcl
terraform {
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
}

provider "kind" {}

provider "kubernetes" {
  host                   = kind_cluster.default.endpoint
  client_certificate     = kind_cluster.default.client_certificate
  client_key             = kind_cluster.default.client_key
  cluster_ca_certificate = kind_cluster.default.cluster_ca_certificate
}

provider "helm" {
  kubernetes {
    host                   = kind_cluster.default.endpoint
    client_certificate     = kind_cluster.default.client_certificate
    client_key             = kind_cluster.default.client_key
    cluster_ca_certificate = kind_cluster.default.cluster_ca_certificate
  }
}
```

### 2.2 k8s.tf - Kind 클러스터

```hcl
resource "kind_cluster" "default" {
  name           = "ingress-gateway-demo"
  wait_for_ready = true

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"
      kubeadm_config_patches = [
        <<-EOF
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
        EOF
      ]
      extra_port_mappings {
        container_port = 80
        host_port      = 80
        protocol       = "TCP"
      }
      extra_port_mappings {
        container_port = 443
        host_port      = 443
        protocol       = "TCP"
      }
    }

    node {
      role = "worker"
    }
  }
}
```

### 2.3 modules/infra/infra.tf - ArgoCD 설치

```hcl
resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "7.8.28"
  namespace  = kubernetes_namespace.argocd.metadata[0].name

  set {
    name  = "configs.params.server.insecure"
    value = "true"
  }

  set {
    name  = "configs.secret.argocdServerAdminPassword"
    value = "$2a$10$rRyBsGSHK6.uc8fntPwVIuLVHgsAhAX7TcdrqW/RBER3kHIwBiP3C" # password
  }
}
```

---

## 3. Helm Charts 구현

### 3.1 echo-server

**Chart.yaml:**
```yaml
apiVersion: v2
name: echo-server
description: Echo Server for Ingress/Gateway demo
version: 0.1.0
appVersion: "latest"
```

**values.yaml:**
```yaml
replicaCount: 1

image:
  repository: kenshin579/echo-server
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 50m
    memory: 64Mi
```

**templates/deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "echo-server.fullname" . }}
  labels:
    {{- include "echo-server.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "echo-server.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "echo-server.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: 80
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /ping
              port: 80
          readinessProbe:
            httpGet:
              path: /ping
              port: 80
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

### 3.2 Ingress - nginx-ingress

**Chart.yaml:**
```yaml
apiVersion: v2
name: nginx-ingress
description: NGINX Ingress Controller wrapper
version: 0.1.0

dependencies:
  - name: ingress-nginx
    version: 4.12.0
    repository: https://kubernetes.github.io/ingress-nginx
```

**values.yaml:**
```yaml
ingress-nginx:
  controller:
    service:
      type: NodePort
    hostPort:
      enabled: true
    nodeSelector:
      ingress-ready: "true"
    tolerations:
      - key: "node-role.kubernetes.io/control-plane"
        operator: "Equal"
        effect: "NoSchedule"
```

### 3.3 Ingress - ingress-routes

**templates/ingress.yaml:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Values.ingress.name }}
  annotations:
    {{- toYaml .Values.ingress.annotations | nindent 4 }}
spec:
  ingressClassName: {{ .Values.ingress.className }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType }}
            backend:
              service:
                name: {{ .serviceName }}
                port:
                  number: {{ .servicePort }}
          {{- end }}
    {{- end }}
```

**values.yaml:**
```yaml
ingress:
  name: echo-server-ingress
  className: nginx
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
  hosts:
    - host: echo.local
      paths:
        - path: /
          pathType: Prefix
          serviceName: echo-server
          servicePort: 80
```

### 3.4 Gateway - gateway-api-crds

**Chart.yaml:**
```yaml
apiVersion: v2
name: gateway-api-crds
description: Gateway API CRDs
version: 0.1.0
appVersion: "1.2.0"
```

**values.yaml:**
```yaml
version: v1.2.0
experimental: true
```

**templates/crds.yaml:**
```yaml
# CRDs are installed via URL in ArgoCD Application
# https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml
```

### 3.5 Gateway - nginx-gateway

**Chart.yaml:**
```yaml
apiVersion: v2
name: nginx-gateway
description: NGINX Gateway Fabric wrapper
version: 0.1.0

dependencies:
  - name: nginx-gateway-fabric
    version: 2.2.1
    repository: oci://ghcr.io/nginx/charts
```

**values.yaml:**
```yaml
nginx-gateway-fabric:
  service:
    type: NodePort
  nginxGateway:
    gwAPIExperimentalFeatures:
      enable: true
```

### 3.6 Gateway - gateway-routes

**templates/gateway.yaml:**
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: {{ .Values.gateway.name }}
  namespace: {{ .Release.Namespace }}
spec:
  gatewayClassName: {{ .Values.gateway.className }}
  listeners:
    {{- range .Values.gateway.listeners }}
    - name: {{ .name }}
      port: {{ .port }}
      protocol: {{ .protocol }}
      allowedRoutes:
        namespaces:
          from: {{ .allowedRoutes.from }}
    {{- end }}
```

**templates/httproutes.yaml:**
```yaml
{{- range .Values.httpRoutes }}
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: {{ .name }}
  namespace: {{ $.Release.Namespace }}
spec:
  parentRefs:
    - name: {{ $.Values.gateway.name }}
      namespace: {{ $.Release.Namespace }}
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

**values.yaml:**
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

httpRoutes:
  - name: echo-server-route
    hostnames:
      - echo.local
    rules:
      - matches:
          - path:
              type: PathPrefix
              value: /
        backendRefs:
          - name: echo-server
            namespace: app
            port: 80
```

---

## 4. ArgoCD Bootstrap 구현

### 4.1 apps.yaml - echo-server 배포

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: echo-server
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/kenshin579/tutorials-go
    targetRevision: HEAD
    path: cloud/ingress-gateway/charts/echo-server
  destination:
    server: https://kubernetes.default.svc
    namespace: app
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

### 4.2 infra-ingress.yaml

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: ingress-infra
  namespace: argocd
spec:
  generators:
    - list:
        elements:
          - name: nginx-ingress
            namespace: ingress-nginx
            path: cloud/ingress-gateway/charts/ingress/nginx-ingress
          - name: ingress-routes
            namespace: app
            path: cloud/ingress-gateway/charts/ingress/ingress-routes
  template:
    metadata:
      name: "{{name}}"
    spec:
      project: default
      source:
        repoURL: https://github.com/kenshin579/tutorials-go
        targetRevision: HEAD
        path: "{{path}}"
      destination:
        server: https://kubernetes.default.svc
        namespace: "{{namespace}}"
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
```

### 4.3 infra-gateway.yaml

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: gateway-infra
  namespace: argocd
spec:
  generators:
    - list:
        elements:
          - name: gateway-api-crds
            namespace: gateway
            path: cloud/ingress-gateway/charts/gateway/gateway-api-crds
          - name: nginx-gateway
            namespace: gateway
            path: cloud/ingress-gateway/charts/gateway/nginx-gateway
          - name: gateway-routes
            namespace: gateway
            path: cloud/ingress-gateway/charts/gateway/gateway-routes
  template:
    metadata:
      name: "{{name}}"
    spec:
      project: default
      source:
        repoURL: https://github.com/kenshin579/tutorials-go
        targetRevision: HEAD
        path: "{{path}}"
      destination:
        server: https://kubernetes.default.svc
        namespace: "{{namespace}}"
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
```

---

## 5. Makefile 구현

```makefile
.PHONY: all tf-init tf-apply tf-destroy clean

TERRAFORM_DIR := terraform

all: tf-init tf-apply

tf-init:
	cd $(TERRAFORM_DIR) && terraform init

tf-apply:
	cd $(TERRAFORM_DIR) && terraform apply -auto-approve

tf-destroy:
	cd $(TERRAFORM_DIR) && terraform destroy -auto-approve

# Ingress 방식 배포
deploy-ingress:
	kubectl apply -f bootstrap/apps.yaml
	kubectl apply -f bootstrap/infra-ingress.yaml

# Gateway 방식 배포
deploy-gateway:
	kubectl apply -f bootstrap/apps.yaml
	kubectl apply -f bootstrap/infra-gateway.yaml

# 테스트
test-ingress:
	@echo "Testing Ingress..."
	curl -s -H "Host: echo.local" http://localhost/ping | jq .

test-gateway:
	@echo "Testing Gateway..."
	curl -s -H "Host: echo.local" http://localhost/ping | jq .

# ArgoCD 접속 정보
argocd-info:
	@echo "ArgoCD URL: http://localhost:8080"
	@echo "Username: admin"
	@echo "Password: password"
	kubectl port-forward svc/argocd-server -n argocd 8080:80 &

clean: tf-destroy
```

---

## 6. 테스트 방법

### 6.1 Ingress 테스트

```bash
# 1. 인프라 생성
make tf-init && make tf-apply

# 2. Ingress 배포
make deploy-ingress

# 3. 배포 확인
kubectl get pods -n ingress-nginx
kubectl get pods -n app
kubectl get ingress -n app

# 4. 테스트
curl -H "Host: echo.local" http://localhost/ping
```

### 6.2 Gateway API 테스트

```bash
# 1. 인프라 생성
make tf-init && make tf-apply

# 2. Gateway 배포
make deploy-gateway

# 3. 배포 확인
kubectl get pods -n gateway
kubectl get pods -n app
kubectl get gateway -n gateway
kubectl get httproute -n gateway

# 4. 테스트
curl -H "Host: echo.local" http://localhost/ping
```
